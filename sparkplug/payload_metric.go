package sparkplug

import (
	"encoding/json"
	"fmt"
	"reflect"

	any "github.com/golang/protobuf/ptypes/any"
)

// all of this insanity to convert datatype to a string when JSON serializing
// this is dumb and i hate it

type payloadPropertyValue struct {
	Type   uint32      `json:"type,omitempty"`
	IsNull bool        `json:"is_null,omitempty"`
	Value  interface{} `json:"value,omitempty"`
}

type payloadPropertySet struct {
	Keys    []string                `json:"keys,omitempty"` // Names of the properties
	Values  []*payloadPropertyValue `json:"values,omitempty"`
	Details []*any.Any              `json:"details,omitempty"`
}

type payloadMetaData struct {
	IsMultiPart bool       `json:"is_multi_part,omitempty"`
	ContentType string     `json:"content_type,omitempty"` // Content/Media type
	Size        uint64     `json:"size,omitempty"`         // File size, String size, Multi-part size, etc
	Seq         uint64     `json:"seq,omitempty"`          // Sequence number for multi-part messages
	FileName    string     `json:"file_name,omitempty"`    // File name
	FileType    string     `json:"file_type,omitempty"`    // File type (i.e. xml, json, txt, cpp, etc)
	Md5         string     `json:"md5,omitempty"`          // md5 of data
	Description string     `json:"description,omitempty"`  // Could be anything such as json or xml of custom properties
	Details     []*any.Any `json:"details,omitempty"`
}

type payloadMetric struct {
	Name         string                 `json:"name,omitempty"`          // Metric name - should only be included on birth
	Alias        uint64                 `json:"alias,omitempty"`         // Metric alias - tied to name on birth and included in all later DATA messages
	Timestamp    uint64                 `json:"timestamp,omitempty"`     // Timestamp associated with data acquisition time
	Datatype     string                 `json:"datatype,omitempty"`      // DataType of the metric/tag value
	IsHistorical bool                   `json:"is_historical,omitempty"` // If this is historical data and should not update real time tag
	IsTransient  bool                   `json:"is_transient,omitempty"`  // Tells consuming clients such as MQTT Engine to not store this as a tag
	IsNull       bool                   `json:"is_null,omitempty"`       // If this is null - explicitly say so rather than using -1, false, etc for some datatypes.
	Metadata     *payloadMetaData       `json:"metadata,omitempty"`      // Metadata for the payload
	Properties   *payloadPropertySet    `json:"properties,omitempty"`
	Value        isPayload_Metric_Value `json:"value,omitempty"`
}

// valueToJsonType will convert a sparkplug datatype to a JSON type,
// one of: number, string, boolean, array
func (x *Payload_Metric) ValueToJsonType() (JsonType, error) {
	if x == nil {
		return nil, fmt.Errorf("nil metric")
	}
	return MetricValueToJsonType(x)
}

func (x *Payload_Metric) MarshalJSON() ([]byte, error) {
	var pMetadata *payloadMetaData
	if x.Metadata != nil {
		pMetadata = &payloadMetaData{
			IsMultiPart: x.Metadata.IsMultiPart,
			ContentType: x.Metadata.ContentType,
			Size:        x.Metadata.Size,
			Seq:         x.Metadata.Seq,
			FileName:    x.Metadata.FileName,
			FileType:    x.Metadata.FileType,
			Md5:         x.Metadata.Md5,
			Description: x.Metadata.Description,
			Details:     x.Metadata.Details,
		}
	}

	var pPropertySet *payloadPropertySet
	if x.Properties != nil {
		ppValues := make([]*payloadPropertyValue, len(x.Properties.Values))
		if x.Properties.Values != nil {
			for j, value := range x.Properties.Values {
				if value == nil {
					continue
				}
				if value.Value == nil {
					continue
				}
				var ppValue interface{}
				valueType := reflect.TypeOf(value.Value)
				switch valueType {
				case reflect.TypeOf(&Payload_PropertyValue_IntValue{}):
					ppValue = value.GetIntValue()
				case reflect.TypeOf(&Payload_PropertyValue_LongValue{}):
					ppValue = value.GetLongValue()
				case reflect.TypeOf(&Payload_PropertyValue_FloatValue{}):
					ppValue = value.GetFloatValue()
				case reflect.TypeOf(&Payload_PropertyValue_DoubleValue{}):
					ppValue = value.GetDoubleValue()
				case reflect.TypeOf(&Payload_PropertyValue_BooleanValue{}):
					ppValue = value.GetBooleanValue()
				case reflect.TypeOf(&Payload_PropertyValue_StringValue{}):
					ppValue = value.GetStringValue()
				case reflect.TypeOf(&Payload_PropertyValue_PropertysetValue{}):
					ppValue = value.GetPropertysetValue()
				case reflect.TypeOf(&Payload_PropertyValue_PropertysetsValue{}):
					ppValue = value.GetPropertysetsValue()
				case reflect.TypeOf(&Payload_PropertyValue_ExtensionValue{}):
					ppValue = value.GetExtensionValue()
				default:
					ppValue = nil
				}
				ppValues[j] = &payloadPropertyValue{
					Type:   value.Type,
					IsNull: value.IsNull,
					Value:  ppValue,
				}
			}
		}
		pPropertySet = &payloadPropertySet{
			Keys:    x.Properties.Keys,
			Values:  ppValues,
			Details: x.Properties.Details,
		}
	}

	metric := &payloadMetric{
		Name:         x.Name,
		Alias:        x.Alias,
		Timestamp:    x.Timestamp,
		Datatype:     DataType_name[int32(x.Datatype)],
		IsHistorical: x.IsHistorical,
		IsTransient:  x.IsTransient,
		IsNull:       x.IsNull,
		Metadata:     pMetadata,
		Properties:   pPropertySet,
		Value:        x.Value,
	}

	return json.Marshal(metric)
}

func (x *Payload_Metric) UnmarshalJSON(data []byte) error {
	var obj map[string]json.RawMessage
	jErr := json.Unmarshal(data, &obj)
	if jErr != nil {
		return jErr
	}

	var metric payloadMetric

	if name, ok := obj["name"]; ok {
		if err := json.Unmarshal(name, &metric.Name); err != nil {
			return err
		}
	}

	if alias, ok := obj["alias"]; ok {
		if err := json.Unmarshal(alias, &metric.Alias); err != nil {
			return err
		}
	}

	if timestamp, ok := obj["timestamp"]; ok {
		if err := json.Unmarshal(timestamp, &metric.Timestamp); err != nil {
			return err
		}
	}

	if datatype, ok := obj["datatype"]; ok {
		if err := json.Unmarshal(datatype, &metric.Datatype); err != nil {
			return err
		}
	}

	if isHistorical, ok := obj["is_historical"]; ok {
		if err := json.Unmarshal(isHistorical, &metric.IsHistorical); err != nil {
			return err
		}
	}

	if isTransient, ok := obj["is_transient"]; ok {
		if err := json.Unmarshal(isTransient, &metric.IsTransient); err != nil {
			return err
		}
	}

	if isNull, ok := obj["is_null"]; ok {
		if err := json.Unmarshal(isNull, &metric.IsNull); err != nil {
			return err
		}
	}

	if metadata, ok := obj["metadata"]; ok {
		var x payloadMetaData
		if err := json.Unmarshal(metadata, &x); err != nil {
			return err
		}
		metric.Metadata = &x
	}

	if properties, ok := obj["properties"]; ok {
		var x payloadPropertySet
		if err := json.Unmarshal(properties, &x); err != nil {
			return err
		}
		metric.Properties = &x
	}

	if value, ok := obj["value"]; ok {
		switch metric.Datatype {
		case "Boolean":
			var x Payload_Metric_BooleanValue
			if err := json.Unmarshal(value, &x); err != nil {
				return err
			}
			metric.Value = &x
		case "Int8":
			fallthrough
		case "Int16":
			fallthrough
		case "Int32":
			fallthrough
		case "Int64":
			fallthrough
		case "UInt8":
			fallthrough
		case "UInt16":
			fallthrough
		case "UInt32":
			var x Payload_Metric_IntValue
			if err := json.Unmarshal(value, &x); err != nil {
				return err
			}
			metric.Value = &x
		case "UInt64":
			var x Payload_Metric_LongValue
			if err := json.Unmarshal(value, &x); err != nil {
				return err
			}
			metric.Value = &x
		case "Float":
			var x Payload_Metric_FloatValue
			if err := json.Unmarshal(value, &x); err != nil {
				return err
			}
			metric.Value = &x
		case "Double":
			var x Payload_Metric_DoubleValue
			if err := json.Unmarshal(value, &x); err != nil {
				return err
			}
			metric.Value = &x
		case "String":
			var x Payload_Metric_StringValue
			if err := json.Unmarshal(value, &x); err != nil {
				return err
			}
			metric.Value = &x
		case "DateTime":
			var x Payload_Metric_LongValue
			if err := json.Unmarshal(value, &x); err != nil {
				return err
			}
			metric.Value = &x
		case "DataSet":
			var x Payload_Metric_DatasetValue
			if err := json.Unmarshal(value, &x); err != nil {
				return err
			}
			metric.Value = &x
		case "UUID":
			var x Payload_Metric_StringValue
			if err := json.Unmarshal(value, &x); err != nil {
				return err
			}
			metric.Value = &x
		default:
			return fmt.Errorf("unknown DataType value: %s", metric.Datatype)
		}
	}

	x.Alias = metric.Alias
	x.Timestamp = metric.Timestamp
	x.Name = metric.Name
	x.Datatype = uint32(DataType_value[metric.Datatype])
	x.IsHistorical = metric.IsHistorical
	x.IsTransient = metric.IsTransient
	x.IsNull = metric.IsNull

	if metric.Metadata != nil {
		x.Metadata = &Payload_MetaData{
			IsMultiPart: metric.Metadata.IsMultiPart,
			ContentType: metric.Metadata.ContentType,
			Size:        metric.Metadata.Size,
			Seq:         metric.Metadata.Seq,
			FileName:    metric.Metadata.FileName,
			FileType:    metric.Metadata.FileType,
			Md5:         metric.Metadata.Md5,
			Description: metric.Metadata.Description,
			Details:     metric.Metadata.Details,
		}
	}

	if metric.Properties != nil {
		propValues := make([]*Payload_PropertyValue, len(metric.Properties.Values))
		for _, value := range metric.Properties.Values {
			switch value.Type {
			case 1:
				propValues = append(propValues, &Payload_PropertyValue{Type: 1, IsNull: value.IsNull, Value: &Payload_PropertyValue_IntValue{IntValue: uint32(value.Value.(int32))}})
			case 2:
				propValues = append(propValues, &Payload_PropertyValue{Type: 2, IsNull: value.IsNull, Value: &Payload_PropertyValue_LongValue{LongValue: uint64(value.Value.(int64))}})
			case 3:
				propValues = append(propValues, &Payload_PropertyValue{Type: 3, IsNull: value.IsNull, Value: &Payload_PropertyValue_FloatValue{FloatValue: float32(value.Value.(float64))}})
			case 4:
				propValues = append(propValues, &Payload_PropertyValue{Type: 4, IsNull: value.IsNull, Value: &Payload_PropertyValue_DoubleValue{DoubleValue: value.Value.(float64)}})
			case 5:
				propValues = append(propValues, &Payload_PropertyValue{Type: 5, IsNull: value.IsNull, Value: &Payload_PropertyValue_BooleanValue{BooleanValue: value.Value.(bool)}})
			case 6:
				propValues = append(propValues, &Payload_PropertyValue{Type: 6, IsNull: value.IsNull, Value: &Payload_PropertyValue_StringValue{StringValue: value.Value.(string)}})
			case 7:
				propValues = append(propValues, &Payload_PropertyValue{Type: 7, IsNull: value.IsNull, Value: &Payload_PropertyValue_PropertysetValue{PropertysetValue: value.Value.(*Payload_PropertySet)}})
			case 8:
				propValues = append(propValues, &Payload_PropertyValue{Type: 8, IsNull: value.IsNull, Value: &Payload_PropertyValue_PropertysetsValue{PropertysetsValue: value.Value.(*Payload_PropertySetList)}})
			case 9:
				propValues = append(propValues, &Payload_PropertyValue{Type: 9, IsNull: value.IsNull, Value: &Payload_PropertyValue_ExtensionValue{ExtensionValue: value.Value.(*Payload_PropertyValue_PropertyValueExtension)}})
			default:
				switch value.Value.(type) {
				case int32:
					propValues = append(propValues, &Payload_PropertyValue{Type: value.Type, IsNull: value.IsNull, Value: &Payload_PropertyValue_IntValue{IntValue: uint32(value.Value.(int32))}})
				case int64:
					propValues = append(propValues, &Payload_PropertyValue{Type: value.Type, IsNull: value.IsNull, Value: &Payload_PropertyValue_LongValue{LongValue: uint64(value.Value.(int64))}})
				case float32:
					propValues = append(propValues, &Payload_PropertyValue{Type: value.Type, IsNull: value.IsNull, Value: &Payload_PropertyValue_FloatValue{FloatValue: float32(value.Value.(float64))}})
				case float64:
					propValues = append(propValues, &Payload_PropertyValue{Type: value.Type, IsNull: value.IsNull, Value: &Payload_PropertyValue_DoubleValue{DoubleValue: value.Value.(float64)}})
				case bool:
					propValues = append(propValues, &Payload_PropertyValue{Type: value.Type, IsNull: value.IsNull, Value: &Payload_PropertyValue_BooleanValue{BooleanValue: value.Value.(bool)}})
				case string:
					propValues = append(propValues, &Payload_PropertyValue{Type: value.Type, IsNull: value.IsNull, Value: &Payload_PropertyValue_StringValue{StringValue: value.Value.(string)}})
				default:
					propValues = append(propValues, &Payload_PropertyValue{Type: value.Type, IsNull: value.IsNull, Value: &Payload_PropertyValue_ExtensionValue{ExtensionValue: value.Value.(*Payload_PropertyValue_PropertyValueExtension)}})
				}

			}
		}

		x.Properties = &Payload_PropertySet{
			Keys:    metric.Properties.Keys,
			Values:  propValues,
			Details: metric.Properties.Details,
		}

	}
	x.Value = metric.Value

	return nil
}
