package service

import (
	"fmt"

	"github.com/american-factory-os/glowplug/json_type"
	"github.com/american-factory-os/glowplug/sparkplug"
)

// valueToJsonType will convert a sparkplug datatype to a JSON type,
// one of: number, string, boolean, array
func PayloadMetricToJsonType(x *sparkplug.Payload_Metric) (json_type.JsonType, error) {
	if x == nil {
		return nil, fmt.Errorf("nil metric")
	}
	return json_type.MetricValueToJsonType(x)
}
