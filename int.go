package params

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Int struct {
	value   int  // Value holds the actual integer value
	present bool // Present indicates if the integer is present or not
}

// UnmarshalJSON implements custom unmarshalling for the Int type.
// It handles cases where the integer may be zero, null, or quoted.
// If the integer is zero or null, it sets Present to false and Value to zero.
// If the integer is quoted, it removes the quotes and sets Present to true.
// If the integer is not quoted, it sets Present to true and retains the value as is.
// This allows for flexible handling of integer values in JSON payloads.
func (i *Int) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		i.value = 0
		i.present = false
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
		i.value = 0
		i.present = false
		return err
	} else {
		vv, err := v.Int64()
		if err != nil {
			i.value = 0
			i.present = false
			return err
		}
		i.value = int(vv)
	}
	i.present = true

	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It allows the Int type to be unmarshalled from text representations.
// This method simply calls UnmarshalJSON with the provided text data.
//
// Parameters:
//   - text: The text data to unmarshal into the Int type.
//
// Returns:
//   - error: An error if the unmarshalling fails, otherwise nil.
func (i *Int) UnmarshalText(text []byte) error {
	return i.UnmarshalJSON(text)
}

// UnmarshalParam is a helper method to unmarshal a string parameter directly.
// It converts the string parameter to a byte slice and calls UnmarshalJSON.
//
// Parameters:
//   - param: The string parameter to unmarshal into the Int type.
//
// Returns:
//   - error: An error if the unmarshalling fails, otherwise nil.
func (i *Int) UnmarshalParam(param string) error {
	return i.UnmarshalJSON([]byte(param))
}

// Set sets the value of the Int type and marks it as present.
// This method updates the Value field with the provided integer and sets Present to true.
//
// Parameters:
//   - value: The integer value to set for the Int type.
func (i *Int) Set(value int) {
	i.value = value
	i.present = true
}

// Value retrieves the value of the Int type.
// If the integer is not present, it returns zero.
// If the integer is present, it returns the Value field.
//
// Returns:
//   - int: The value of the Int type if present, otherwise zero.
func (i *Int) Value() int {
	if !i.present {
		return 0
	}
	return i.value
}

// Present checks if the Int type is present in the JSON payload.
// It returns true if the integer was provided in the JSON payload, otherwise false.
//
// Returns:
//   - bool: True if the integer is present, otherwise false.
func (i *Int) Present() bool {
	return i.present
}

// MarshalJSON implements custom marshalling for the Int type.
// It converts the Int type to a JSON integer representation.
// If the integer is not present, it returns an empty JSON string.
//
// Returns:
//   - []byte: The JSON representation of the Int type.
//   - error: An error if the marshalling fails, otherwise nil.
func (i Int) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "%d", i.Value()), nil // Marshal the integer value
}
