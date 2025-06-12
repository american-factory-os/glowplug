package service

import (
	"github.com/american-factory-os/glowplug/sparkplug"
)

// WebsocketMetricMessage represents a JSON SparkplugB metric sent over websocket
type WebsocketMetricMessage struct {
	Topic     *sparkplug.Topic   `json:"topic"`
	Alias     uint64             `json:"alias"`
	Name      string             `json:"name"`
	Value     sparkplug.JsonType `json:"value"`
	Timestamp uint64             `json:"timestamp"`
}
