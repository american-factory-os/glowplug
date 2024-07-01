package sparkplug

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestDataTypeMarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		dt      DataType
		want    []byte
		wantErr bool
	}{
		{
			name: "Test data type string",
			dt:   DataType_String,
			want: []byte("\"String\""),
		},
		{
			name: "Test data type int8",
			dt:   DataType_Int8,
			want: []byte("\"Int8\""),
		},
		{
			name: "Test data type int16",
			dt:   DataType_Int16,
			want: []byte("\"Int16\""),
		},
		{
			name: "Test data type int32",
			dt:   DataType_Int32,
			want: []byte("\"Int32\""),
		},
		{
			name: "Test data type int64",
			dt:   DataType_Int64,
			want: []byte("\"Int64\""),
		},
		{
			name: "Test data type uint8",
			dt:   DataType_UInt8,
			want: []byte("\"UInt8\""),
		},
		{
			name: "Test data type uint16",
			dt:   DataType_UInt16,
			want: []byte("\"UInt16\""),
		},
		{
			name: "Test data type uint32",
			dt:   DataType_UInt32,
			want: []byte("\"UInt32\""),
		},
		{
			name: "Test data type uint64",
			dt:   DataType_UInt64,
			want: []byte("\"UInt64\""),
		},
		{
			name: "Test data type float",
			dt:   DataType_Float,
			want: []byte("\"Float\""),
		},
		{
			name: "Test data type double",
			dt:   DataType_Double,
			want: []byte("\"Double\""),
		},
		{
			name: "Test data type boolean",
			dt:   DataType_Boolean,
			want: []byte("\"Boolean\""),
		},
		{
			name: "Test data type string",
			dt:   DataType_String,
			want: []byte("\"String\""),
		},
		{
			name: "Test data type datetime",
			dt:   DataType_DateTime,
			want: []byte("\"DateTime\""),
		},
		{
			name: "Test data type text",
			dt:   DataType_Text,
			want: []byte("\"Text\""),
		},
		{
			name: "Test data type uuid",
			dt:   DataType_UUID,
			want: []byte("\"UUID\""),
		},
		{
			name: "Test data type dataset",
			dt:   DataType_DataSet,
			want: []byte("\"DataSet\""),
		},
		{
			name: "Test data type bytes",
			dt:   DataType_Bytes,
			want: []byte("\"Bytes\""),
		},
		{
			name: "Test data type file",
			dt:   DataType_File,
			want: []byte("\"File\""),
		},
		{
			name: "Test data type template",
			dt:   DataType_Template,
			want: []byte("\"Template\""),
		},
		{
			name: "Test data type property set",
			dt:   DataType_PropertySet,
			want: []byte("\"PropertySet\""),
		},
		{
			name: "Test data type property set list",
			dt:   DataType_PropertySetList,
			want: []byte("\"PropertySetList\""),
		},
		{
			name: "Test data type int8 array",
			dt:   DataType_Int8Array,
			want: []byte("\"Int8Array\""),
		},
		{
			name: "Test data type int16 array",
			dt:   DataType_Int16Array,
			want: []byte("\"Int16Array\""),
		},
		{
			name: "Test data type int32 array",
			dt:   DataType_Int32Array,
			want: []byte("\"Int32Array\""),
		},
		{
			name: "Test data type int64 array",
			dt:   DataType_Int64Array,
			want: []byte("\"Int64Array\""),
		},
		{
			name: "Test data type uint8 array",
			dt:   DataType_UInt8Array,
			want: []byte("\"UInt8Array\""),
		},
		{
			name: "Test data type uint16 array",
			dt:   DataType_UInt16Array,
			want: []byte("\"UInt16Array\""),
		},
		{
			name: "Test data type uint32 array",
			dt:   DataType_UInt32Array,
			want: []byte("\"UInt32Array\""),
		},
		{
			name: "Test data type uint64 array",
			dt:   DataType_UInt64Array,
			want: []byte("\"UInt64Array\""),
		},
		{
			name: "Test data type float array",
			dt:   DataType_FloatArray,
			want: []byte("\"FloatArray\""),
		},
		{
			name: "Test data type double array",
			dt:   DataType_DoubleArray,
			want: []byte("\"DoubleArray\""),
		},
		{
			name: "Test data type boolean array",
			dt:   DataType_BooleanArray,
			want: []byte("\"BooleanArray\""),
		},
		{
			name: "Test data type string array",
			dt:   DataType_StringArray,
			want: []byte("\"StringArray\""),
		},
		{
			name: "Test data type datetime array",
			dt:   DataType_DateTimeArray,
			want: []byte("\"DateTimeArray\""),
		},
		{
			name: "Test data type Unknown",
			dt:   DataType_Unknown,
			want: []byte("\"Unknown\""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.dt)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestDataTypeMarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("TestDataTypeMarshalJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestDataTypeUnmarshalJSON(t *testing.T) {
	type testFormat struct {
		name    string
		dtBytes []byte
		want    DataType
		wantErr bool
	}

	tests := []testFormat{}

	for dataTypeNum, dataTypeStr := range DataType_name {
		test := testFormat{
			name:    string(dataTypeStr),
			dtBytes: []byte("\"" + string(dataTypeStr) + "\""),
			want:    DataType(dataTypeNum),
			wantErr: false,
		}
		tests = append(tests, test)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dt DataType
			err := json.Unmarshal(tt.dtBytes, &dt)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestDataTypeUnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if dt != tt.want {
				t.Errorf("TestDataTypeUnmarshalJSON() = %v, want %v", dt, tt.want)
			}
		})
	}
}
