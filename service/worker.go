package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync/atomic"

	"github.com/american-factory-os/glowplug/sparkplug"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
)

const statReportInterval = 1000

const (
	STATE_STOPPED uint32 = 0
	STATE_RUNNING uint32 = 1
)

type Result struct {
	err         error
	sourceTopic string
	payload     *sparkplug.Payload
	topic       *sparkplug.Topic
}

type Worker interface {
	Run() error
	AddMessage(msg Message) error
	Stop()
	Capacity() (current int, size int)
}

type worker struct {
	state         *atomic.Uint32
	logger        *log.Logger
	size          int
	messages      chan Message
	results       chan Result
	rdb           *redis.UniversalClient
	publishBroker *mqtt.Client
	total         uint64
	errors        uint64
	seen          map[string]bool
}

func (w *worker) Stop() {
	w.state.Store(STATE_STOPPED)
	close(w.messages)
	if w.rdb != nil {
		rdb := *w.rdb
		rdb.Close()
	}
	if publishBroker, err := w.getPublishBroker(); err == nil {
		defer publishBroker.Disconnect(250)
	}
}

func (w *worker) processResult(result Result) error {

	if result.err != nil {
		if strings.Contains(result.err.Error(), "invalid wire-format data") {
			return fmt.Errorf("sparkplug %w from topic %s", result.err, result.sourceTopic)
		}
		return fmt.Errorf("error processing message, %w", result.err)
	}

	if result.payload == nil {
		return fmt.Errorf("no payload found")
	}

	if result.payload.Metrics == nil {
		return nil
	}

	if len(result.payload.Metrics) == 0 {
		return nil
	}

	switch result.topic.Command {
	case sparkplug.NBIRTH:
		fallthrough
	case sparkplug.NDATA:
		fallthrough
	case sparkplug.DBIRTH:
		fallthrough
	case sparkplug.DDATA:
		for _, metric := range result.payload.Metrics {
			key := keyFromSparkplugMetric(*result.topic, metric)

			var value interface{} = nil

			switch metric.Datatype {
			case uint32(sparkplug.DataType_Float):
				value = metric.GetFloatValue()
			case uint32(sparkplug.DataType_Double):
				value = metric.GetDoubleValue()
			case uint32(sparkplug.DataType_Int8):
				fallthrough
			case uint32(sparkplug.DataType_Int16):
				fallthrough
			case uint32(sparkplug.DataType_Int32):
				fallthrough
			case uint32(sparkplug.DataType_Int64):
				fallthrough
			case uint32(sparkplug.DataType_UInt8):
				fallthrough
			case uint32(sparkplug.DataType_UInt16):
				fallthrough
			case uint32(sparkplug.DataType_UInt32):
				fallthrough
			case uint32(sparkplug.DataType_UInt64):
				value = metric.GetIntValue()
			case uint32(sparkplug.DataType_Boolean):
				value = metric.GetBooleanValue()
			case uint32(sparkplug.DataType_String):
				value = metric.GetStringValue()
			case uint32(sparkplug.DataType_DateTime):
				// Date time value as uint64 value representing milliseconds since epoch (Jan 1, 1970)
				value = metric.GetIntValue()
			case uint32(sparkplug.DataType_Text):
				value = metric.GetStringValue()
			case uint32(sparkplug.DataType_UUID):
				value = metric.GetStringValue()
			case uint32(sparkplug.DataType_DataSet):
				value = metric.GetStringValue()
			case uint32(sparkplug.DataType_Bytes):
				value = metric.GetBytesValue()
			case uint32(sparkplug.DataType_File):
				value = metric.GetBytesValue()
			case uint32(sparkplug.DataType_Template):
				value = metric.GetBytesValue()
			case uint32(sparkplug.DataType_PropertySet):
				fallthrough
			case uint32(sparkplug.DataType_PropertySetList):
				fallthrough
			case uint32(sparkplug.DataType_Unknown):
				fallthrough
			case uint32(sparkplug.DataType_Int8Array):
				fallthrough
			case uint32(sparkplug.DataType_Int16Array):
				fallthrough
			case uint32(sparkplug.DataType_Int32Array):
				fallthrough
			case uint32(sparkplug.DataType_Int64Array):
				fallthrough
			case uint32(sparkplug.DataType_UInt8Array):
				fallthrough
			case uint32(sparkplug.DataType_UInt16Array):
				fallthrough
			case uint32(sparkplug.DataType_UInt32Array):
				fallthrough
			case uint32(sparkplug.DataType_UInt64Array):
				fallthrough
			case uint32(sparkplug.DataType_FloatArray):
				fallthrough
			case uint32(sparkplug.DataType_DoubleArray):
				fallthrough
			case uint32(sparkplug.DataType_BooleanArray):
				fallthrough
			case uint32(sparkplug.DataType_StringArray):
				fallthrough
			case uint32(sparkplug.DataType_DateTimeArray):
				// All array types use the bytes_value field of the Metric value field. They are simply little-endian packed byte arrays.
				value = metric.GetBytesValue()
			default:
				w.logger.Println("unknown datatype value", metric.Datatype)
			}

			// json encode the value
			if bytes, err := json.Marshal(value); err != nil {
				return err
			} else {
				value = string(bytes)
			}

			// report new metric seen
			typeName := sparkplug.DataType_name[int32(metric.Datatype)]
			_, seen := w.seen[key]
			if !seen {
				w.seen[key] = true
				w.logger.Println("new metric encountered", key, typeName, value)
			}

			// pipeline redis commands
			if w.rdb != nil {
				rdb := *w.rdb
				if cmds, err := rdb.Pipelined(context.TODO(), func(pipeliner redis.Pipeliner) error {

					if !seen {
						// save human readable metric type in a redis hash
						pipeliner.HSet(context.TODO(), HASH_METRIC_TYPES, key, typeName)
					}

					// store the metric value in a redis set
					pipeliner.Set(context.TODO(), key, value, 0)

					// publish metric value to redis channel
					pipeliner.Publish(context.TODO(), key, value)

					return nil
				}); err != nil {
					return err
				} else {
					for _, cmd := range cmds {
						if cmd.Err() != nil && cmd.Err() != redis.Nil {
							return fmt.Errorf("redis cmd error %w", cmd.Err())
						}
					}
				}
			}

			if w.publishBroker != nil {
				// publish metric value to mqtt
				go func(t *sparkplug.Topic, m *sparkplug.Payload_Metric) {
					topic := topicFromSparkplugMetric(*t, m)
					if publishBroker, err := w.getPublishBroker(); err == nil {
						if token := publishBroker.Publish(topic, 0, false, value); token.Wait() && token.Error() != nil {
							log.Println("unable to publish to mqtt", token.Error())
						}
					}

				}(result.topic, metric)
			}

		}
	}

	return nil
}

func (w *worker) processResults() error {

	for {
		result, ok := <-w.results

		if !ok {
			return fmt.Errorf("channel closed, no longer processing results")
		}

		if w.state.Load() == STATE_STOPPED {
			return fmt.Errorf("worker pool stopped, dropping result on topic %s", result.sourceTopic)
		}

		err := w.processResult(result)
		if err != nil {
			w.logger.Println(err)
			w.errors++
		} else {
			w.total++
		}

		if w.total > 0 && w.total%statReportInterval == 0 {
			w.logger.Printf("processed %d messages, %d errors\n", w.total, w.errors)
		}
	}
}

func (w *worker) getPublishBroker() (mqtt.Client, error) {
	if w.publishBroker == nil {
		return nil, errors.New("publish broker not available")
	}

	return *w.publishBroker, nil
}

func (w *worker) Run() error {
	w.state.Store(STATE_RUNNING)

	go w.processResults()

	for {
		msg, ok := <-w.messages

		if !ok {
			return fmt.Errorf("channel closed, no longer processing messages")
		}

		if w.state.Load() == STATE_STOPPED {
			return fmt.Errorf("worker pool stopped, dropping message on topic %s", msg.topic)
		}

		topic, tErr := sparkplug.ToTopic(msg.topic)
		if tErr != nil {
			w.results <- Result{
				sourceTopic: msg.topic,
				err:         tErr,
			}
			continue
		}

		var processCmd bool

		// process node and device birth and data commands
		switch topic.Command {
		case sparkplug.NBIRTH:
			fallthrough
		case sparkplug.NDATA:
			fallthrough
		case sparkplug.DBIRTH:
			fallthrough
		case sparkplug.DDATA:
			processCmd = true
		default:
			processCmd = false
		}

		if processCmd {
			var payload sparkplug.Payload
			err := proto.Unmarshal(msg.payload, &payload)
			if err != nil {
				w.results <- Result{
					sourceTopic: msg.topic,
					err:         err,
				}
				continue
			}

			w.results <- Result{
				err:         nil,
				sourceTopic: msg.topic,
				payload:     &payload,
				topic:       topic,
			}
		}
	}
}

func (w *worker) AddMessage(msg Message) error {
	if w.state.Load() == STATE_STOPPED {
		return errors.New("worker pool stopped")
	}

	w.messages <- msg
	return nil
}

// Capacity returns current message capacity and size
func (w *worker) Capacity() (current int, size int) {
	current = cap(w.messages)
	size = w.size
	return
}

func NewWorker(logger *log.Logger, rdb *redis.UniversalClient, publishBroker *mqtt.Client) (Worker, error) {

	size := runtime.NumCPU() * 100

	state := atomic.Uint32{}
	state.Store(STATE_STOPPED)

	return &worker{
		state:         &state,
		logger:        logger,
		size:          size,
		messages:      make(chan Message, size),
		results:       make(chan Result, size),
		rdb:           rdb,
		publishBroker: publishBroker,
		seen:          make(map[string]bool),
	}, nil
}
