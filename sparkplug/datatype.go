package sparkplug

import (
	"encoding/json"
	"fmt"
)

// MarshalJSON marshals a DataType
func (d DataType) MarshalJSON() ([]byte, error) {
	str, ok := DataType_name[int32(d)]
	if !ok {
		return nil, fmt.Errorf("unknown DataType value: %d", d)
	}
	return json.Marshal(str)
}

// MarshalBinary marshals a DataType
func (d DataType) MarshalBinary() ([]byte, error) {
	return d.MarshalJSON()
}

// UnmarshalJSON unmarshals a DataType
func (d *DataType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	for k, v := range DataType_name {
		if v == str {
			*d = DataType(k)
			return nil
		}
	}
	return fmt.Errorf("unknown DataType value: %s", str)
}

// Uint32 returns the DataType as a uint32
func (d DataType) Uint32() uint32 {
	return uint32(d)
}
