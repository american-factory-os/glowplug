package service

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/american-factory-os/glowplug/sparkplug"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	topicDelimiter = "/"
	topicPrefix    = "glowplug"
)

type Message struct {
	topic   string
	payload []byte
}

// The URL format should be scheme://host:port Where "scheme" is one of:
// "mqtt", "tcp", "ssl", or "ws", "host" is the ip-address (or hostname)
// and "port" is the port on which the broker is accepting connections
// Default values for hostname is "127.0.0.1", for schema is "tcp://".
// An example broker URI would look like: tcp://foobar.com:1883
func validateBrokerURI(rawURL string) (bool, error) {
	if rawURL == "" {
		return false, fmt.Errorf("broker URL is empty")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return false, err
	}

	if u.Scheme != "tcp" && u.Scheme != "ssl" && u.Scheme != "ws" && u.Scheme != "mqtt" {
		return false, fmt.Errorf("invalid scheme: %s", u.Scheme)
	}

	if len(u.Host)-len(u.Port()) < 4 {
		return false, fmt.Errorf("invalid host")
	}

	if u.Port() == "" {
		return false, fmt.Errorf("invalid port")
	}

	return true, nil
}

// normalizeTopicName ensures mqtt topics keys are in a standard format
func normalizeTopicName(ns string) string {
	return strings.ReplaceAll(ns, ":", "/")
}

// topicFromSparkplugMetric returns a topic for a given SparkplugB metric
func topicFromSparkplugMetric(topic sparkplug.Topic, metric *sparkplug.Payload_Metric) string {

	var b strings.Builder
	b.WriteString(topicPrefix)
	b.WriteString(topicDelimiter)

	b.WriteString(topic.GroupId)
	b.WriteString(topicDelimiter)
	b.WriteString(topic.EdgeNodeId)
	if topic.HasDevice {
		b.WriteString(topicDelimiter)
		b.WriteString(topic.DeviceId)
	}

	b.WriteString(topicDelimiter)
	b.WriteString(normalizeTopicName(metric.Name))

	return normalizeTopicName(b.String())
}

// brokerClientFromURL returns a mqtt.Client from a given URL
func brokerClientFromURL(rawURL string, handler *mqtt.MessageHandler) (mqtt.Client, error) {

	if _, err := validateBrokerURI(rawURL); err != nil {
		return nil, err
	}

	randClientID := fmt.Sprintf("glowplug-%d", time.Now().UnixNano())
	mqttOpts := mqtt.NewClientOptions().AddBroker(rawURL).SetClientID(randClientID)
	mqttOpts.SetKeepAlive(2 * time.Second)
	mqttOpts.SetPingTimeout(1 * time.Second)

	if handler != nil {
		mqttOpts.SetDefaultPublishHandler(*handler)
	}

	broker := mqtt.NewClient(mqttOpts)
	if token := broker.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return broker, nil
}
