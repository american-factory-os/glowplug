package json_type

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/american-factory-os/glowplug/sparkplug"
	"github.com/gopcua/opcua/ua"
)

var ErrPayloadMetricNil = fmt.Errorf("payload metrics is nil")
var ErrPayloadMetricNilHasProperties = fmt.Errorf("payload metrics is nil, properties are not nil")

// JsonType is an interface that represents either a number, string, boolean, or an array of those types.
// It's intented to serialize single values or arrays of values to JSON
type JsonType interface {
	MarshalJSON() ([]byte, error)
	MarshalBinary() ([]byte, error)
	String() string
	Bytes() []byte
}

type jsonNull struct{}

func (x *jsonNull) MarshalJSON() ([]byte, error) {
	return json.Marshal(nil)
}

func (x *jsonNull) MarshalBinary() ([]byte, error) {
	return json.Marshal(nil)
}

func (x *jsonNull) String() string {
	return ""
}

func (x *jsonNull) Bytes() []byte {
	return nil
}

type jsonArray struct {
	a []interface{}
}

func (x *jsonArray) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.a)
}

func (x *jsonArray) MarshalBinary() ([]byte, error) {
	return json.Marshal(x.a)
}

func (x *jsonArray) String() string {
	return string(x.Bytes())
}

func (x *jsonArray) Bytes() []byte {
	if x == nil {
		panic("nil json array")
	}
	b, e := json.Marshal(x.a)
	if e != nil {
		panic(e)
	}
	return b
}

type jsonNumber struct {
	n interface{}
}

func (x *jsonNumber) MarshalJSON() ([]byte, error) {
	if x.n != nil {
		return json.Marshal(x.n)
	}
	return nil, fmt.Errorf("no value to marshal")
}

func (x *jsonNumber) MarshalBinary() ([]byte, error) {
	if x.n != nil {
		return json.Marshal(x.n)
	}
	return nil, fmt.Errorf("no value to marshal")
}

func (x *jsonNumber) String() string {
	if x.Bytes() != nil {
		return string(x.Bytes())
	}
	return ""
}

func (x *jsonNumber) Bytes() []byte {
	if x == nil {
		return nil
	}

	// in the future, think about how to marshal the number without using the json package which uses exponents, ex:
	// var buf [8]byte
	// binary.BigEndian.PutUint64(buf[:], math.Float64bits(f))

	b, err := json.Marshal(x.n)
	if err != nil {
		panic(err)
	}

	return b

}

type jsonString struct {
	s *string
}

func (x *jsonString) MarshalJSON() ([]byte, error) {
	if x.s != nil {
		return json.Marshal(*x.s)
	}
	return nil, fmt.Errorf("no value to marshal")
}

func (x *jsonString) MarshalBinary() ([]byte, error) {
	if x.s != nil {
		return json.Marshal(*x.s)
	}
	return nil, fmt.Errorf("no value to marshal")
}

func (x *jsonString) String() string {
	if x.Bytes() != nil {
		return string(x.Bytes())
	}
	return ""
}

func (x *jsonString) Bytes() []byte {
	if x.s != nil {
		return []byte(*x.s)
	}
	return nil
}

type jsonBool struct {
	b *bool
}

func (x *jsonBool) MarshalJSON() ([]byte, error) {
	if x.b != nil {
		return json.Marshal(*x.b)
	}
	return nil, fmt.Errorf("no value to marshal")
}

func (x *jsonBool) MarshalBinary() ([]byte, error) {
	if x.b != nil {
		return json.Marshal(*x.b)
	}
	return nil, fmt.Errorf("no value to marshal")
}

func (x *jsonBool) String() string {
	if x.b != nil {
		return string(x.Bytes())
	}
	return "false"
}

func (x *jsonBool) Bytes() []byte {
	if x.b != nil {
		if *x.b {
			return []byte("true")
		}
	}
	return []byte("false")
}

func newJsonNumber[T uint64 | uint32 | uint16 | uint8 | uint | int64 | int32 | int16 | int8 | int | float64 | float32](v T) JsonType {
	return &jsonNumber{
		n: v,
	}
}

func newJsonString(s string) JsonType {
	return &jsonString{
		s: &s,
	}
}

func newJsonBool(b bool) JsonType {
	return &jsonBool{
		b: &b,
	}
}

// newJsonArray only supports arrays of basic types
func newJsonArray[T uint64 | uint32 | uint16 | uint8 | uint | int64 | int32 | int16 | int8 | int | float64 | float32 | string | bool](a []T) (JsonType, error) {
	if a == nil {
		return nil, nil
	}

	data := make([]interface{}, len(a))

	for i, v := range a {
		data[i] = v
	}

	return &jsonArray{a: data}, nil
}

func newJsonNull() JsonType {
	return &jsonNull{}
}

// MetricValueToJsonType will convert a sparkplug datatype to a JSON type,
// one of: number, string, boolean, array
func MetricValueToJsonType(metric *sparkplug.Payload_Metric) (JsonType, error) {

	if metric == nil {
		return nil, fmt.Errorf("metric is nil, can't convert to JsonType")
	}

	if metric.Value == nil {
		return newJsonNull(), nil
	}

	// cast to int32 because we know the datatype is valid per proto
	datatype := int32(metric.Datatype)

	name, ok := sparkplug.DataType_name[datatype]
	if !ok {
		name = ""
	}

	switch name {
	case "Int8":
		return newJsonNumber(int8(metric.GetIntValue())), nil
	case "Int16":
		return newJsonNumber(int16(metric.GetIntValue())), nil
	case "Int32":
		return newJsonNumber(int32(metric.GetIntValue())), nil
	case "Int64":
		return newJsonNumber(int64(metric.GetIntValue())), nil
	case "UInt8":
		fallthrough
	case "UInt16":
		fallthrough
	case "UInt32":
		fallthrough
	case "UInt64":
		return newJsonNumber(uint64(metric.GetLongValue())), nil
	case "Float":
		return newJsonNumber(float32(metric.GetFloatValue())), nil
	case "Double":
		return newJsonNumber(float64(metric.GetDoubleValue())), nil
	case "Boolean":
		return newJsonBool(metric.GetBooleanValue()), nil
	case "DateTime":
		return newJsonNumber(int64(metric.GetIntValue())), nil
	case "String":
		fallthrough
	case "Text":
		fallthrough
	case "UUID":
		fallthrough
	case "DataSet":
		return newJsonString(metric.GetStringValue()), nil
	case "Bytes":
		fallthrough
	case "File":
		return newJsonString(string(metric.GetBytesValue())), nil
	case "Template":
		fallthrough
	case "PropertySet":
		fallthrough
	case "PropertySetList":
		fallthrough
	case "Int8Array":
		fallthrough
	case "Int16Array":
		fallthrough
	case "Int32Array":
		fallthrough
	case "Int64Array":
		fallthrough
	case "UInt8Array":
		fallthrough
	case "UInt16Array":
		fallthrough
	case "UInt32Array":
		fallthrough
	case "UInt64Array":
		fallthrough
	case "FloatArray":
		fallthrough
	case "DoubleArray":
		fallthrough
	case "BooleanArray":
		fallthrough
	case "StringArray":
		fallthrough
	case "DateTimeArray":
		fallthrough
	case "Unknown":
		fallthrough
	default:
		return nil, fmt.Errorf("sparkplug datatype %d %s is currently unsupported", datatype, name)
	}

}

// NodeValueToJsonType will convert a OPC UA Node type to a JSON type,
// one of: number, string, boolean, array
func NodeValueToJsonType(variant *ua.Variant) (JsonType, error) {

	if variant == nil {
		return nil, fmt.Errorf("nil variant")
	}

	if variant.Value() == nil {
		return nil, fmt.Errorf("nil value")
	}

	datatype := variant.Type()

	fail := func(v *ua.Variant) (JsonType, error) {
		return nil, fmt.Errorf("ua datatype %d is currently unsupported, reflected type is %v", v.Type(), reflect.TypeOf(v.Value()))
	}

	switch datatype {
	case ua.TypeIDNull:
		return fail(variant)
	case ua.TypeIDBoolean:
		return newJsonBool(variant.Bool()), nil
	case ua.TypeIDSByte:
		return fail(variant)
	case ua.TypeIDByte:
		return fail(variant)
	case ua.TypeIDInt16:
		fallthrough
	case ua.TypeIDInt32:
		fallthrough
	case ua.TypeIDInt64:
		return newJsonNumber(int64(variant.Int())), nil
	case ua.TypeIDUint16:
		fallthrough
	case ua.TypeIDUint32:
		fallthrough
	case ua.TypeIDUint64:
		return newJsonNumber(uint64(variant.Uint())), nil
	case ua.TypeIDFloat:
		return fail(variant)
	case ua.TypeIDDouble:
		return newJsonNumber(variant.Float()), nil
	case ua.TypeIDString:
		return newJsonString(variant.String()), nil
	case ua.TypeIDDateTime:
		return fail(variant)
	case ua.TypeIDGUID:
		return newJsonString(variant.GUID().String()), nil
	case ua.TypeIDByteString:
		return fail(variant)
	case ua.TypeIDXMLElement:
		return fail(variant)
	case ua.TypeIDNodeID:
		return fail(variant)
	case ua.TypeIDExpandedNodeID:
		return fail(variant)
	case ua.TypeIDStatusCode:
		return fail(variant)
	case ua.TypeIDQualifiedName:
		return fail(variant)
	case ua.TypeIDLocalizedText:
		return newJsonString(variant.LocalizedText().Text), nil
	case ua.TypeIDExtensionObject:
		return fail(variant)
	case ua.TypeIDDataValue:
		return fail(variant)
	case ua.TypeIDVariant:
		return fail(variant)
	case ua.TypeIDDiagnosticInfo:
		return fail(variant)
	default:
		return nil, fmt.Errorf("ua datatype %d is currently unsupported, reflected type is %v", datatype, reflect.TypeOf(variant.Value()))
	}

}
