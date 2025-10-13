package params

import (
	"encoding/json"
)

// Structure for handling strings in JSON payloads
// This structure allows for the presence of a string to be explicitly indicated,
type String struct {
	value   string // The actual string value
	present bool   // Indicates if the string is present in the JSON payload
}

// UnmarshalJSON implements custom unmarshalling for the String type.
// It handles cases where the string may be empty, null, or quoted.
// If the string is empty or null, it sets Present to false and Value to an empty string.
// If the string is quoted, it removes the quotes and sets Present to true.
// If the string is not quoted, it sets Present to true and retains the value as is.
// This allows for flexible handling of string values in JSON payloads.
//
// Parameters:
//   - data: The JSON data to unmarshal into the String type.
//
// Returns:
//   - error: An error if the unmarshalling fails, otherwise nil.
func (s *String) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		s.value = ""
		s.present = false
		return nil
	}

	if err := json.Unmarshal(data, &s.value); err != nil {
		s.value = ""
		s.present = false
		return err
	}
	s.present = true

	return nil
}

// Set sets the value of the String type and marks it as present.
// This method updates the Value field with the provided string and sets Present to true.
//
// Parameters:
//   - value: The string value to set for the String type.
func (s *String) Set(value string) {
	s.value = value
	s.present = true
}

// MarshalJSON implements custom marshalling for the String type.
// It converts the String type to a JSON string representation.
// If the string is not present, it returns an empty JSON string.
// If the string is present, it returns the value wrapped in quotes.
//
// Returns:
//   - []byte: The JSON representation of the String type.
//   - error: An error if the marshalling fails, otherwise nil.
func (s String) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Value())
}

// GetJSON returns the JSON representation of the String type.
// It marshals the Value field into a JSON string.
// If the marshaling fails, it returns an empty string.
//
// Returns:
//   - string: The JSON representation of the String type, or an empty string if marshaling fails.
func (s *String) GetJSON() string {
	b, err := json.Marshal(s.value)
	if err != nil {
		return ""
	}
	return string(b)
}

// Present checks if the String type is present in the JSON payload.
// It returns true if the string was provided in the JSON payload, otherwise false.
//
// Returns:
//   - bool: True if the string is present, otherwise false.
func (s *String) Present() bool {
	return s.present
}

// Value retrieves the actual string value of the String type.
// If the string is not present, it returns an empty string.
// If the string is present, it returns the Value field.
//
// Returns:
//   - string: The actual string value if present, otherwise an empty string.
func (s *String) Value() string {
	if !s.present {
		return ""
	}
	return s.value
}
