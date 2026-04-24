package params

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Structure for handling UUIDs in JSON payloads
// This structure allows for the presence of a UUID to be explicitly indicated,
type UUID struct {
	value   uuid.UUID // The actual UUID value
	present bool      // Indicates if the UUID is present in the JSON payload
}

// UnmarshalJSON implements custom unmarshalling for the UUID type.
// It handles cases where the UUID may be empty, null, or quoted.
// If the UUID is empty or null, it sets Present to false and Value to an empty UUID.
// If the UUID is quoted, it removes the quotes and sets Present to true.
// If the UUID is not quoted, it sets Present to true and retains the value as is.
// This allows for flexible handling of UUID values in JSON payloads.
//
// Parameters:
//   - data: The JSON data to unmarshal into the UUID type.
//
// Returns:
//   - error: An error if the unmarshalling fails, otherwise nil.
func (u *UUID) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		u.value = uuid.UUID{}
		u.present = false
		return nil
	}

	if err := json.Unmarshal(data, &u.value); err != nil {
		u.value = uuid.UUID{}
		u.present = false
		return err
	}
	u.present = true

	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It allows the UUID type to be unmarshalled from text representations.
// This method simply calls UnmarshalJSON with the provided text data.
//
// Parameters:
//   - text: The text data to unmarshal into the UUID type.
//
// Returns:
//   - error: An error if the unmarshalling fails, otherwise nil.
func (u *UUID) UnmarshalText(text []byte) error {
	return u.UnmarshalJSON(text)
}

// UnmarshalParam implements the custom parameter unmarshalling for the UUID type.
// It allows the UUID type to be unmarshalled directly from a string parameter.
// This method simply calls UnmarshalJSON with the provided string data.
//
// Parameters:
//   - param: The string parameter to unmarshal into the UUID type.
//
// Returns:
//   - error: An error if the unmarshalling fails, otherwise nil.
func (u *UUID) UnmarshalParam(param string) error {
	return u.UnmarshalJSON([]byte(param))
}

// Set sets the value of the UUID type and marks it as present.
// This method updates the Value field with the provided UUID and sets Present to true.
//
// Parameters:
//   - value: The UUID value to set for the UUID type.
func (u *UUID) Set(value uuid.UUID) {
	u.value = value
	u.present = true
}

// MarshalJSON implements custom marshalling for the UUID type.
// It converts the UUID type to a JSON string representation.
// If the UUID is not present, it returns an empty JSON string.
// If the UUID is present, it returns the value wrapped in quotes.
//
// Returns:
//   - []byte: The JSON representation of the UUID type.
//   - error: An error if the marshalling fails, otherwise nil.
func (u UUID) MarshalJSON() ([]byte, error) {
	if !u.present {
		return []byte(`""`), nil
	}
	return json.Marshal(u.Value())
}

// GetJSON returns the JSON representation of the UUID type.
// It marshals the Value field into a JSON string.
// If the marshaling fails, it returns an empty string.
//
// Returns:
//   - string: The JSON representation of the UUID type, or an empty string if marshaling fails.
func (u *UUID) GetJSON() string {
	b, err := json.Marshal(u.value)
	if err != nil {
		return ""
	}
	return string(b)
}

// Present checks if the UUID type is present in the JSON payload.
// It returns true if the UUID was provided in the JSON payload, otherwise false.
//
// Returns:
//   - bool: True if the UUID is present, otherwise false.
func (u *UUID) Present() bool {
	return u.present
}

// Value retrieves the actual UUID value of the UUID type.
// If the UUID is not present, it returns an empty UUID.
// If the UUID is present, it returns the Value field.
//
// Returns:
//   - uuid.UUID: The actual UUID value if present, otherwise an empty UUID.
func (u *UUID) Value() uuid.UUID {
	if !u.present {
		return uuid.UUID{}
	}
	return u.value
}
