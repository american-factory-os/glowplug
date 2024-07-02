package sparkplug

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewJsonNumber(t *testing.T) {
	intValue := int64(42)
	uintValue := uint64(42)
	float32Value := float32(42.0)
	float64Value := float64(42.0)

	jtInt := newJsonInt64(intValue)
	assert.Equal(t, fmt.Sprint(intValue), jtInt.String())
	assert.Equal(t, []byte(fmt.Sprint(intValue)), jtInt.Bytes())

	jtIntBytes, intErr := jtInt.MarshalJSON()
	assert.Nil(t, intErr)
	assert.Equal(t, fmt.Sprint(intValue), string(jtIntBytes))

	jtUint := newJsonUInt64(uintValue)
	assert.Equal(t, fmt.Sprint(uintValue), jtUint.String())
	assert.Equal(t, []byte(fmt.Sprint(uintValue)), jtUint.Bytes())

	jtUintBytes, uintErr := jtUint.MarshalJSON()
	assert.Nil(t, uintErr)
	assert.Equal(t, fmt.Sprint(uintValue), string(jtUintBytes))

	jtFloat32 := newJsonFloat32(float32Value)
	assert.Equal(t, fmt.Sprint(float32Value), jtFloat32.String())
	assert.Equal(t, []byte(fmt.Sprint(float32Value)), jtFloat32.Bytes())

	jtFloat32Bytes, float32Err := jtFloat32.MarshalJSON()
	assert.Nil(t, float32Err)
	assert.Equal(t, fmt.Sprint(float32Value), string(jtFloat32Bytes))

	jtFloat64 := newJsonFloat64(float64Value)
	assert.Equal(t, fmt.Sprint(float64Value), jtFloat64.String())
	assert.Equal(t, []byte(fmt.Sprint(float64Value)), jtFloat64.Bytes())

	jtFloat64Bytes, float64Err := jtFloat64.MarshalJSON()
	assert.Nil(t, float64Err)
	assert.Equal(t, fmt.Sprint(float64Value), string(jtFloat64Bytes))
	assert.Equal(t, []byte(fmt.Sprint(float64Value)), jtFloat64.Bytes())
}

func TestNewJsonString(t *testing.T) {
	str := "test string"
	jt := newJsonString(str)
	assert.Equal(t, str, jt.String())
	assert.Equal(t, []byte(str), jt.Bytes())

	jtBytes, err := jt.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("\"%s\"", str), string(jtBytes))
}

func TestNewJsonBool(t *testing.T) {
	bTrue := true
	jt := newJsonBool(bTrue)
	assert.Equal(t, fmt.Sprint(bTrue), jt.String())
	assert.Equal(t, []byte("true"), jt.Bytes())

	jtBytes, err := jt.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprint(bTrue), string(jtBytes))

	bFalse := false
	jtFalse := newJsonBool(bFalse)
	assert.Equal(t, fmt.Sprint(bFalse), jtFalse.String())
	assert.Equal(t, []byte("false"), jtFalse.Bytes())

	jtFalseBytes, falseErr := jtFalse.MarshalJSON()
	assert.Nil(t, falseErr)
	assert.Equal(t, fmt.Sprint(bFalse), string(jtFalseBytes))
}
