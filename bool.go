package params

import (
	"fmt"
	"strings"
)

type Bool struct {
	value   bool // Value holds the actual boolean value
	present bool // Present indicates if the boolean is present or not
}

// UnmarshalJSON implements custom unmarshalling for the Bool type.
// It handles cases where the boolean may be true, false, null, or quoted.
// If the boolean is null, it sets Present to false and Value to false.
// If the boolean is quoted, it removes the quotes and sets Present to true.
// If the boolean is not quoted, it sets Present to true and retains the value as is.
// This allows for flexible handling of boolean values in JSON payloads.
func (b *Bool) UnmarshalJSON(data []byte) error {
	b.present = false

	if len(data) == 0 || string(data) == "null" {
		b.value = false
		return nil
	}

	str := strings.ToLower(strings.Trim(string(data), `"`))

	switch str {
	case "true":
		b.value = true
	case "false":
		b.value = false
	default:
		return fmt.Errorf("invalid boolean format: %s", string(data))
	}

	b.present = true

	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It allows the Bool type to be unmarshalled from text representations.
// This method simply calls UnmarshalJSON with the provided text data.
//
// Parameters:
//   - text: The text data to unmarshal into the Bool type.
//
// Returns:
//   - error: An error if the unmarshalling fails, otherwise nil.
func (b *Bool) UnmarshalText(text []byte) error {
	return b.UnmarshalJSON(text)
}

// UnmarshalParam is a helper method to unmarshal a string parameter directly.
// It converts the string parameter to a byte slice and calls UnmarshalJSON.
//
// Parameters:
//   - param: The string parameter to unmarshal into the Bool type.
//
// Returns:
//   - error: An error if the unmarshalling fails, otherwise nil.
func (b *Bool) UnmarshalParam(param string) error {
	return b.UnmarshalJSON([]byte(param))
}

// Set sets the value of the Bool type and marks it as present.
// This method updates the Value field with the provided boolean and sets Present to true.
//
// Parameters:
//   - value: The boolean value to set for the Bool type.
func (b *Bool) Set(value bool) {
	b.value = value
	b.present = true
}

// Value retrieves the value of the Bool type.
// If the boolean is not present, it returns false.
// If the boolean is present, it returns the Value field.
//
// Returns:
//   - bool: The value of the Bool type if present, otherwise false.
func (b *Bool) Value() bool {
	if !b.present {
		return false
	}
	return b.value
}

// Present checks if the Bool type is present in the JSON payload.
// It returns true if the boolean was provided in the JSON payload, otherwise false.
//
// Returns:
//   - bool: True if the boolean is present, otherwise false.
func (b *Bool) Present() bool {
	return b.present
}

// MarshalJSON implements custom marshalling for the Bool type.
// It converts the Bool type to a JSON boolean representation.
// If the boolean is not present, it returns an empty JSON string.
//
// Returns:
//   - []byte: The JSON representation of the Bool type.
//   - error: An error if the marshalling fails, otherwise nil.
func (b Bool) MarshalJSON() ([]byte, error) {
	if !b.present {
		return []byte("null"), nil
	}
	if b.value {
		return []byte("true"), nil
	}
	return []byte("false"), nil
}
