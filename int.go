package params

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Int struct {
	Value   int  `json:"-"`
	Present bool `json:"-"` // Present indicates if the integer is present or not
}

// UnmarshalJSON implements custom unmarshalling for the Int type.
// It handles cases where the integer may be zero, null, or quoted.
// If the integer is zero or null, it sets Present to false and Value to zero.
// If the integer is quoted, it removes the quotes and sets Present to true.
// If the integer is not quoted, it sets Present to true and retains the value as is.
// This allows for flexible handling of integer values in JSON payloads.
func (i *Int) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		i.Value = 0
		i.Present = false
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
		i.Value = 0
		i.Present = false
		return err
	} else {
		vv, err := v.Int64()
		if err != nil {
			i.Value = 0
			i.Present = false
			return err
		}
		i.Value = int(vv)
	}
	i.Present = true

	return nil
}

// Set sets the value of the Int type and marks it as present.
// This method updates the Value field with the provided integer and sets Present to true.
//
// Parameters:
//   - value: The integer value to set for the Int type.
func (i *Int) Set(value int) {
	i.Value = value
	i.Present = true
}

// Get retrieves the value of the Int type.
// If the integer is not present, it returns zero.
// If the integer is present, it returns the Value field.
//
// Returns:
//   - int: The value of the Int type if present, otherwise zero.
func (i *Int) Get() int {
	if !i.Present {
		return 0
	}
	return i.Value
}

// MarshalJSON implements custom marshalling for the Int type.
// It converts the Int type to a JSON integer representation.
// If the integer is not present, it returns an empty JSON string.
//
// Returns:
//   - []byte: The JSON representation of the Int type.
//   - error: An error if the marshalling fails, otherwise nil.
func (i Int) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "%d", i.Get()), nil // Marshal the integer value
}
