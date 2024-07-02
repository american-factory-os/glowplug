package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/american-factory-os/glowplug/sparkplug"
	"github.com/gopcua/opcua/ua"
	"github.com/redis/go-redis/v9"
)

const (
	HASH_METRIC_TYPES = "glowplug:metric_types"
)

const (
	keyDelimiter      = ":"
	keyPrefixGlowplug = "glowplug"
	keyPrefixOpcua    = "opcua"
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
	b.WriteString(keyPrefixGlowplug)
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

// keyFromUaNodeId returns a namespace for a given OPCUA Node ID
func keyFromUaNodeId(productURI string, nodeID *ua.NodeID) (string, error) {

	if nodeID == nil {
		return "", fmt.Errorf("nodeID is nil")
	}

	var b strings.Builder
	b.WriteString(keyPrefixGlowplug)
	b.WriteString(keyDelimiter)

	b.WriteString(keyPrefixOpcua)
	b.WriteString(keyDelimiter)

	b.WriteString(productURI)
	b.WriteString(keyDelimiter)

	b.WriteString(fmt.Sprint(nodeID.Namespace()))
	b.WriteString(keyDelimiter)

	switch nodeID.Type() {
	case ua.NodeIDTypeTwoByte:
		fallthrough
	case ua.NodeIDTypeFourByte:
		fallthrough
	case ua.NodeIDTypeNumeric:
		b.WriteString("i")
		b.WriteString(keyDelimiter)
		b.WriteString(fmt.Sprint(nodeID.IntID()))
	case ua.NodeIDTypeString:
		fallthrough
	case ua.NodeIDTypeGUID:
		fallthrough
	case ua.NodeIDTypeByteString:
		b.WriteString("s")
		b.WriteString(keyDelimiter)
		b.WriteString(nodeID.StringID())
	default:
		return "", fmt.Errorf("unsupported nodeID type: %s", nodeID.Type().String())
	}

	return normalizeKey(b.String()), nil
}
