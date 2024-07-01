package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/american-factory-os/glowplug/sparkplug"
	"github.com/redis/go-redis/v9"
)

const (
	HASH_METRIC_TYPES = "glowplug:metric_types"
)

const (
	keyDelimiter = ":"
	keyPrefix    = "glowplug"
)

// normalizeKey ensures redis keys are in a standard format
func normalizeKey(ns string) string {
	ns = strings.ReplaceAll(ns, " ", "_")
	ns = strings.ReplaceAll(ns, "/", ":")
	return strings.ToLower(ns)
}

// keyFromSparkplugMetric returns a namespace for a given SparkplugB metric
func keyFromSparkplugMetric(topic sparkplug.Topic, metric *sparkplug.Payload_Metric) string {

	var b strings.Builder
	b.WriteString(keyPrefix)
	b.WriteString(keyDelimiter)

	b.WriteString(topic.GroupId)
	b.WriteString(keyDelimiter)
	b.WriteString(topic.EdgeNodeId)
	if topic.HasDevice {
		b.WriteString(keyDelimiter)
		b.WriteString(topic.DeviceId)
	}

	b.WriteString(keyDelimiter)
	b.WriteString(normalizeKey(metric.Name))

	return normalizeKey(b.String())
}

func NewRedis(url string) (*redis.UniversalClient, error) {

	redisOpts, urlErr := redis.ParseURL(url)
	if urlErr != nil {
		return nil, fmt.Errorf("unable to parse redis URL, %w", urlErr)
	}

	opts := redis.UniversalOptions{
		Addrs: []string{redisOpts.Addr},
	}

	rdb := redis.NewUniversalClient(&opts)

	err := rdb.Ping(context.TODO()).Err()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to redis, %w", err)
	}

	return &rdb, nil
}
