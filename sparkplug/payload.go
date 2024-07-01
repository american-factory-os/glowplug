package sparkplug

import (
	"encoding/json"
	"fmt"
)

func (x *Payload_DataSet_DataSetValue) UnmarshalJSON(bytes []byte) error {

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}

	if value, ok := raw["Value"]; !ok {
		return fmt.Errorf("value not found")
	} else {

		var valueRaw map[string]json.RawMessage
		if err := json.Unmarshal(value, &valueRaw); err != nil {
			return err
		}

		if _, ok := valueRaw["IntValue"]; ok {
			var val Payload_DataSet_DataSetValue_IntValue
			if err := json.Unmarshal(value, &val); err != nil {
				return err
			}
			x.Value = &val
			return nil
		}

		if _, ok := valueRaw["LongValue"]; ok {
			var val Payload_DataSet_DataSetValue_LongValue
			if err := json.Unmarshal(value, &val); err != nil {
				return err
			}
			x.Value = &val
			return nil
		}

		if _, ok := valueRaw["FloatValue"]; ok {
			var val Payload_DataSet_DataSetValue_FloatValue
			if err := json.Unmarshal(value, &val); err != nil {
				return err
			}
			x.Value = &val
			return nil
		}

		if _, ok := valueRaw["DoubleValue"]; ok {
			var val Payload_DataSet_DataSetValue_DoubleValue
			if err := json.Unmarshal(value, &val); err != nil {
				return err
			}
			x.Value = &val
			return nil
		}

		if _, ok := valueRaw["BooleanValue"]; ok {
			var val Payload_DataSet_DataSetValue_BooleanValue
			if err := json.Unmarshal(value, &val); err != nil {
				return err
			}
			x.Value = &val
			return nil
		}

		if _, ok := valueRaw["StringValue"]; ok {
			var val Payload_DataSet_DataSetValue_StringValue
			if err := json.Unmarshal(value, &val); err != nil {
				return err
			}
			x.Value = &val
			return nil
		}

		if _, ok := valueRaw["ExtensionValue"]; ok {
			var val Payload_DataSet_DataSetValue_ExtensionValue
			if err := json.Unmarshal(value, &val); err != nil {
				return err
			}
			x.Value = &val
			return nil
		}

		if _, ok := valueRaw["DataSetValueExtension"]; ok {
			return fmt.Errorf("DataSetValueExtension not implemented")
		}

		return fmt.Errorf("value not found in Value field %v", valueRaw)
	}
}
