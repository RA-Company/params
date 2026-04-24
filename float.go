package params

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Float struct {
	value   float64 // Value holds the actual float value
	present bool    // Present indicates if the float is present or not
}

// UnmarshalJSON implements custom unmarshalling for the Float type.
// It handles cases where the float may be zero, null, or quoted.
// If the float is zero or null, it sets Present to false and Value to zero.
// If the float is quoted, it removes the quotes and sets Present to true.
// If the float is not quoted, it sets Present to true and retains the value as is.
// This allows for flexible handling of float values in JSON payloads.
func (f *Float) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		f.value = 0
		f.present = false
		return nil
	}

	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.UseNumber()

	//err := json.Unmarshal(data, &alt)
	//if err != nil {
	//return err
	//}

	var v json.Number

	if err := json.Unmarshal(data, &v); err != nil {
		f.value = 0
		f.present = false
		return err
	} else {
		vv, err := v.Float64()
		if err != nil {
			f.value = 0
			f.present = false
			return err
		}
		f.value = vv
	}
	f.present = true

	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It allows the Float type to be unmarshalled from text representations.
// This method simply calls UnmarshalJSON with the provided text data.
//
// Parameters:
//   - text: The text data to unmarshal into the Float type.
//
// Returns:
//   - error: An error if the unmarshalling fails, otherwise nil.
func (f *Float) UnmarshalText(text []byte) error {
	return f.UnmarshalJSON(text)
}

// UnmarshalParam is a helper method to unmarshal a string parameter directly.
// It converts the string parameter to a byte slice and calls UnmarshalJSON.
//
// Parameters:
//   - param: The string parameter to unmarshal into the Float type.
//
// Returns:
//   - error: An error if the unmarshalling fails, otherwise nil.
func (f *Float) UnmarshalParam(param string) error {
	return f.UnmarshalJSON([]byte(param))
}

// Set sets the value of the Float type and marks it as present.
// This method updates the Value field with the provided float and sets Present to true.
//
// Parameters:
//   - value: The float value to set for the Float type.
func (f *Float) Set(value float64) {
	f.value = value
	f.present = true
}

// Value retrieves the value of the Float type.
// If the float is not present, it returns zero.
// If the float is present, it returns the Value field.
//
// Returns:
//   - float64: The value of the Float type if present, otherwise zero.
func (f *Float) Value() float64 {
	if !f.present {
		return 0
	}
	return f.value
}

// Present checks if the Float type is present in the JSON payload.
// It returns true if the float was provided in the JSON payload, otherwise false.
//
// Returns:
//   - bool: True if the float is present, otherwise false.
func (f *Float) Present() bool {
	return f.present
}

// MarshalJSON implements custom marshalling for the Float type.
// It converts the Float type to a JSON float representation.
// If the float is not present, it returns an empty JSON string.
//
// Returns:
//   - []byte: The JSON representation of the Float type.
//   - error: An error if the marshalling fails, otherwise nil.
func (f Float) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "%f", f.Value()), nil // Marshal the float value
}
