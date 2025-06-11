package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/american-factory-os/glowplug/embed"
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
	Run(httpListenAddr string) error
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
	seen          sync.Map
	wss           WebsocketServer
	httpStop      chan bool
}

func (w *worker) Stop() {
	w.state.Store(STATE_STOPPED)

	// signal the http server to stop
	w.httpStop <- true

	// close the Message channel to stop processing
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

	// process each metric in the payload
	for _, metric := range result.payload.Metrics {
		if len(metric.Name) == 0 {
			return fmt.Errorf("empty metric name")
		}

		// convert sparkplug datatype to json type
		jsonType, err := PayloadMetricToJsonType(metric)
		if err != nil {
			return err
		}

		// redis key for the metric
		key := keyFromSparkplugMetric(*result.topic, metric)

		// report new metric seen
		typeName := sparkplug.DataType_name[int32(metric.Datatype)]
		_, seen := w.seen.Load(key)
		if !seen {
			w.seen.Store(key, true)
			w.logger.Printf("first seen: [%s] %s alias:%d %s:%s\n", result.sourceTopic, metric.Name, metric.Alias, typeName, jsonType)
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
				pipeliner.Set(context.TODO(), key, jsonType, 0)

				// publish metric value to redis channel
				pipeliner.Publish(context.TODO(), key, jsonType)

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
			go func(topic *sparkplug.Topic, metric *sparkplug.Payload_Metric, worker *worker) {
				metricTopic := topicFromSparkplugMetric(*topic, metric)
				if publishBroker, err := worker.getPublishBroker(); err == nil {
					if token := publishBroker.Publish(metricTopic, 0, false, jsonType.Bytes()); token.Wait() && token.Error() != nil {
						log.Println("unable to publish to mqtt", metricTopic, token.Error())
					}
				}

			}(result.topic, metric, w)
		}

		// push data to websocket server
		if w.wss.IsRunning() {

			w.wss.PushData(WebsocketMetricMessage{
				Topic: result.topic,
				Alias: metric.GetAlias(),
				Name:  metric.GetName(),
				Value: jsonType,
			})
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

func (w *worker) Run(httpListenAddr string) error {
	w.state.Store(STATE_RUNNING)

	go func() {
		if len(httpListenAddr) > 0 {
			// Create an HTTP server
			server := &http.Server{
				Addr:    httpListenAddr,
				Handler: nil, // Use default handler
			}

			// Serve index.html for the root path ("/")
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				// Serve index.html for exactly "/" or "/index.html"
				if r.URL.Path == "/" || r.URL.Path == "/index.html" {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					if _, err := w.Write([]byte(embed.GetIndexHTML())); err != nil {
						log.Printf("error serving index.html: %v", err)
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					}
					return
				}
				// Return 404 for other paths (except WebSocket)
				http.NotFound(w, r)
			})

			// Register WebSocket handler for "/ws"
			http.Handle("/ws", w.wss)

			// Start server in a goroutine
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("http server error: %v", err)
				}
			}()

			// Log server started
			w.logger.Println("http server started on", httpListenAddr)

			// Wait for a signal
			<-w.httpStop
			w.logger.Println("received shutdown signal, stopping http server...")

			// Create a context with a timeout for graceful shutdown
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Perform graceful shutdown
			if err := server.Shutdown(ctx); err != nil {
				log.Printf("http shutdown error: %v", err)
			}
			w.logger.Println("http server stopped")
		}
	}()

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
	size = w.size
	current = size - len(w.messages)
	if current < 0 {
		current = 0
	}
	return
}

func NewWorker(logger *log.Logger, rdb *redis.UniversalClient, publishBroker *mqtt.Client, wss WebsocketServer) (Worker, error) {

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
		seen:          sync.Map{},
		wss:           wss,
		httpStop:      make(chan bool, 1),
	}, nil
}
