package service

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/american-factory-os/glowplug/sparkplug"
)

// JsonType is an interface that represents either a number, string, boolean, or array
type JsonType interface {
	MarshalJSON() ([]byte, error)
	MarshalBinary() ([]byte, error)
	String() string
	Bytes() []byte
}

type jsonNumber struct {
	i64 *int64
	u64 *uint64
	f32 *float32
	f64 *float64
}

func (x *jsonNumber) MarshalJSON() ([]byte, error) {
	if x.i64 != nil {
		return json.Marshal(*x.i64)
	}
	if x.u64 != nil {
		return json.Marshal(*x.u64)
	}
	if x.f32 != nil {
		return json.Marshal(*x.f32)
	}
	if x.f64 != nil {
		return json.Marshal(*x.f64)
	}
	return nil, fmt.Errorf("no value to marshal")
}

func (x *jsonNumber) MarshalBinary() ([]byte, error) {
	if x.i64 != nil {
		return json.Marshal(*x.i64)
	}
	if x.u64 != nil {
		return json.Marshal(*x.u64)
	}
	if x.f32 != nil {
		return json.Marshal(*x.f32)
	}
	if x.f64 != nil {
		return json.Marshal(*x.f64)
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
	if x.i64 != nil {
		return []byte(strconv.FormatInt(*x.i64, 10))
	}
	if x.u64 != nil {
		return []byte(strconv.FormatUint(*x.u64, 10))
	}
	if x.f32 != nil {
		return []byte(strconv.FormatFloat(float64(*x.f32), 'f', -1, 32))
	}
	if x.f64 != nil {
		return []byte(strconv.FormatFloat(*x.f64, 'f', -1, 64))
	}
	return nil
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

func newJsonInt64(v int64) JsonType {
	return &jsonNumber{
		i64: &v,
	}
}

func newJsonUInt64(v uint64) JsonType {
	return &jsonNumber{
		u64: &v,
	}
}

func newJsonFloat32(v float32) JsonType {
	return &jsonNumber{
		f32: &v,
	}
}

func newJsonFloat64(v float64) JsonType {
	return &jsonNumber{
		f64: &v,
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

// CoerceSparkplugDatatype will convert a sparkplug datatype to a JSON type,
// one of: number, string, boolean, array
func CoerceSparkplugDatatype(datatype uint32, metric *sparkplug.Payload_Metric) (JsonType, error) {

	// cast to int32 because we know the datatype is valid per sparkplug.proto
	name, ok := sparkplug.DataType_name[int32(datatype)]
	if !ok {
		return nil, fmt.Errorf("unknown sparkplug datatype %d", datatype)
	}

	switch name {
	case "Int8":
		fallthrough
	case "Int16":
		fallthrough
	case "Int32":
		fallthrough
	case "Int64":
		return newJsonInt64(int64(metric.GetIntValue())), nil
	case "UInt8":
		fallthrough
	case "UInt16":
		fallthrough
	case "UInt32":
		fallthrough
	case "UInt64":
		return newJsonUInt64(metric.GetLongValue()), nil
	case "Float":
		return newJsonFloat32(metric.GetFloatValue()), nil
	case "Double":
		return newJsonFloat64(metric.GetDoubleValue()), nil
	case "Boolean":
		return newJsonBool(metric.GetBooleanValue()), nil
	case "DateTime":
		return newJsonInt64(int64(metric.GetIntValue())), nil
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
		return nil, fmt.Errorf("unsupported sparkplug datatype %d", datatype)
	}

}
