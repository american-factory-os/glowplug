package json_type

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewJsonNumber(t *testing.T) {
	int8Value := int8(0)
	int64Value := int64(math.MaxInt64)
	uint64Value := uint64(math.MaxUint64)
	float32Value := float32(math.MaxFloat32)
	float64Value := float64(math.MaxFloat64)

	t.Run("int8", func(t *testing.T) {
		jtInt8 := newJsonNumber(int8Value)
		assert.Equal(t, fmt.Sprint(int8Value), jtInt8.String())
		assert.Equal(t, []byte(fmt.Sprint(int8Value)), jtInt8.Bytes())
		jtInt8Bytes, intErr := jtInt8.MarshalJSON()
		assert.Nil(t, intErr)
		assert.Equal(t, fmt.Sprint(int8Value), string(jtInt8Bytes))
	})

	t.Run("int64", func(t *testing.T) {
		jtInt64 := newJsonNumber(int64Value)
		assert.Equal(t, fmt.Sprint(int64Value), jtInt64.String())
		assert.Equal(t, []byte(fmt.Sprint(int64Value)), jtInt64.Bytes())
		jtIntBytes, intErr := jtInt64.MarshalJSON()
		assert.Nil(t, intErr)
		assert.Equal(t, fmt.Sprint(int64Value), string(jtIntBytes))
	})

	t.Run("uint64", func(t *testing.T) {
		jtUint64 := newJsonNumber(uint64Value)
		assert.Equal(t, fmt.Sprint(uint64Value), jtUint64.String())
		assert.Equal(t, []byte(fmt.Sprint(uint64Value)), jtUint64.Bytes())
		jtUintBytes, uintErr := jtUint64.MarshalJSON()
		assert.Nil(t, uintErr)
		assert.Equal(t, fmt.Sprint(uint64Value), string(jtUintBytes))
	})

	t.Run("float32", func(t *testing.T) {
		expectedStr := "3.4028235e+38"
		jtFloat32 := newJsonNumber(float32Value)
		assert.Equal(t, expectedStr, jtFloat32.String())
		assert.Equal(t, []byte(expectedStr), jtFloat32.Bytes())
		jtFloat32Bytes, float32Err := jtFloat32.MarshalJSON()
		assert.Nil(t, float32Err)
		assert.Equal(t, expectedStr, string(jtFloat32Bytes))
	})

	t.Run("float64", func(t *testing.T) {
		expectedStr := "1.7976931348623157e+308"
		jtFloat64 := newJsonNumber(float64Value)
		assert.Equal(t, expectedStr, jtFloat64.String())
		assert.Equal(t, []byte(expectedStr), jtFloat64.Bytes())
		jtFloat64Bytes, float64Err := jtFloat64.MarshalJSON()
		assert.Nil(t, float64Err)
		assert.Equal(t, expectedStr, string(jtFloat64Bytes))
	})
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

func TestNewJsonArray(t *testing.T) {

	t.Run("int64 array", func(t *testing.T) {
		intArray := []int64{1, 2, 3}
		expectedBytes, _ := json.Marshal(intArray)

		jt, err := newJsonArray(intArray)
		assert.Nil(t, err)
		assert.Equal(t, string(expectedBytes), jt.String())
		assert.Equal(t, expectedBytes, jt.Bytes())

		jtBytes, err := jt.MarshalJSON()
		assert.Nil(t, err)
		assert.JSONEq(t, string(expectedBytes), string(jtBytes))
	})

	t.Run("string array", func(t *testing.T) {
		strArray := []string{"a", "b", "c"}
		expectedBytes, _ := json.Marshal(strArray)

		jt, err := newJsonArray(strArray)
		assert.Nil(t, err)
		assert.Equal(t, string(expectedBytes), jt.String())
		assert.Equal(t, expectedBytes, jt.Bytes())

		jtBytes, err := jt.MarshalJSON()
		assert.Nil(t, err)
		assert.JSONEq(t, string(expectedBytes), string(jtBytes))
	})

	t.Run("bool array", func(t *testing.T) {
		boolArray := []bool{true, false, true}
		expectedBytes, _ := json.Marshal(boolArray)

		jt, err := newJsonArray(boolArray)
		assert.Nil(t, err)
		assert.Equal(t, string(expectedBytes), jt.String())
		assert.Equal(t, expectedBytes, jt.Bytes())

		jtBytes, err := jt.MarshalJSON()
		assert.Nil(t, err)
		assert.JSONEq(t, string(expectedBytes), string(jtBytes))
	})
}

func TestNewJsonNil(t *testing.T) {
	jt := newJsonNull()
	assert.Equal(t, "", jt.String())
	assert.Nil(t, jt.Bytes())

	jtBytes, err := jt.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, "null", string(jtBytes))
}
