package main

import (
	"log"
	"os"
	"time"

	"github.com/american-factory-os/glowplug/sparkplug"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"google.golang.org/protobuf/proto"
)

func uint64TimeNow() uint64 {
	return uint64(time.Now().UTC().UnixNano())
}

func floatPayloadMetric(name string, value float32, seq uint64) *sparkplug.Payload {
	return &sparkplug.Payload{
		Seq:       seq,
		Timestamp: uint64TimeNow(),
		Metrics: []*sparkplug.Payload_Metric{
			{
				Name: name,
				Value: &sparkplug.Payload_Metric_FloatValue{
					FloatValue: value,
				},
				Datatype: sparkplug.DataType_Float.Uint32(),
			},
		},
	}
}

func main() {
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("localhost-test")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// our example sparkplug message payload
	groupId := "Plant1:Area3:Line4:Cell2"
	edgeNodeId := "Heater"
	deviceId := "TempSensor"
	value := float32(98.6)

	// device birth
	birthTopic := sparkplug.DeviceBirthTopic(groupId, edgeNodeId, deviceId)
	dBirthPayload := floatPayloadMetric("Current/Celsius", value, 0)
	birthBytes, err := proto.Marshal(dBirthPayload)
	if err != nil {
		panic(err)
	}
	token := c.Publish(birthTopic, 0, false, birthBytes)
	token.Wait()

	// device data
	seq := uint64(0)
	for i := 0; i < 5; i++ {
		value = value + 0.1

		dataTopic := sparkplug.DeviceDataTopic(groupId, edgeNodeId, deviceId)
		dDataPayload := floatPayloadMetric("Current/Celsius", value, seq)
		bytes, err := proto.Marshal(dDataPayload)
		if err != nil {
			panic(err)
		}
		token := c.Publish(dataTopic, 0, false, bytes)
		token.Wait()

		seq = sparkplug.NextSequenceNumber(seq)

		// sleeping 1s between data points...
		time.Sleep(1 * time.Second)
	}

	// sleeping 1s before device death...
	time.Sleep(1 * time.Second)

	// publish device death
	deathTopic := sparkplug.DeviceDeathTopic(groupId, edgeNodeId, deviceId)
	dDeathPayload := sparkplug.Payload{
		Seq:       uint64(1),
		Timestamp: uint64TimeNow(),
	}
	deathBytes, err := proto.Marshal(&dDeathPayload)
	if err != nil {
		panic(err)
	}

	token = c.Publish(deathTopic, 0, false, deathBytes)
	token.Wait()

	defer c.Disconnect(250)
}
