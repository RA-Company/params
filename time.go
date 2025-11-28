package params

import (
	"fmt"
	"strings"
	"time"
)

var timeLayouts = []string{
	time.RFC3339,              // 2025-09-09T13:20:25Z или с оффсетом
	"2006-01-02T15:04:05 MST", // 2025-09-09T13:20:25 UTC
	"2006-01-02 15:04:05",     // 2025-09-09 13:20:25
	"2006-01-02T15:04:05",     // 2025-09-09T13:20:25
}

// Time is a wrapper around time. Time that supports null values and multiple JSON formats.
type Time struct {
	value   time.Time // Value holds the actual time value
	present bool      // Present indicates if the time is present or not
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It supports multiple time formats and null values.
//
// Parameters:
//   - data: JSON data to unmarshal.
//
// Returns:
//   - error: An error if unmarshaling fails, otherwise nil.
func (dst *Time) UnmarshalJSON(data []byte) error {
	dst.value = time.Time{}
	if len(data) == 0 || string(data) == "null" {
		dst.present = false
		return nil
	}

	if string(data) == `""` {
		dst.present = true
		return nil
	}

	dst.present = true

	for _, layout := range timeLayouts {
		if t, err := time.Parse(layout, strings.Trim(string(data), `"`)); err == nil {
			dst.value = t
			return nil
		}
	}

	return fmt.Errorf("invalid time format: %s", string(data))
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It allows the Time type to be unmarshalled from text representations.
// This method simply calls UnmarshalJSON with the provided text data.
//
// Parameters:
//   - text: The text data to unmarshal into the Time type.
//
// Returns:
//   - error: An error if the unmarshalling fails, otherwise nil.
func (dst *Time) UnmarshalText(text []byte) error {
	return dst.UnmarshalJSON(text)
}

// UnmarshalParam implements the custom parameter unmarshalling for the Time type.
// It allows the Time type to be unmarshalled directly from a string parameter.
// This method simply calls UnmarshalJSON with the provided string data.
//
// Parameters:
//   - param: The string parameter to unmarshal into the Time type.
//
// Returns:
//   - error: An error if the unmarshalling fails, otherwise nil.
func (dst *Time) UnmarshalParam(param string) error {
	return dst.UnmarshalJSON([]byte(param))
}

// MarshalJSON implements the json.Marshaler interface.
// It returns "null" if the time is not present.
//
// Returns:
//   - []byte: JSON representation of the time.
//   - error: An error if marshaling fails, otherwise nil.
func (dst *Time) MarshalJSON() ([]byte, error) {
	if !dst.present {
		return []byte("null"), nil
	}
	return dst.value.MarshalJSON()
}

// IsZero checks if the Time is zero or not present.
//
// Returns:
//   - bool: true if the Time is zero or not present, otherwise false.
func (dst *Time) IsZero() bool {
	return !dst.present || dst.value.IsZero()
}

// Format formats the Time using the provided layout.
// If the Time is not present, it returns an empty string.
//
// Parameters:
//   - layout: The layout to format the time.
//
// Returns:
//   - string: Formatted time string or empty string if not present.
func (dst *Time) Format(layout string) string {
	if !dst.present {
		return ""
	}
	return dst.value.Format(layout)
}

// Equal checks if two Time instances are equal.
//
// Returns:
//   - bool: true if both Time instances are equal, otherwise false.
func (dst *Time) String() string {
	if !dst.present {
		return "null"
	}
	return dst.value.String()
}

// Set sets the value of the Time and marks it as present.
// This method updates the Value field with the provided time and sets Present to true.
//
// Parameters:
//   - value: The time value to set for the Time type.
func (dst *Time) Set(value time.Time) {
	dst.value = value
	dst.present = true
}

// Present checks if the Time type is present in the JSON payload.
// It returns true if the time was provided in the JSON payload, otherwise false.
//
// Returns:
//   - bool: True if the time is present, otherwise false.
func (dst *Time) Present() bool {
	return dst.present
}

// Value retrieves the value of the Time type.
// If the time is not present, it returns the zero value of time.Time.
// If the time is present, it returns the Value field.
//
// Returns:
//   - time.Time: The value of the Time type if present, otherwise the zero value of time.Time.
func (dst *Time) Value() time.Time {
	if !dst.present {
		return time.Time{}
	}
	return dst.value
}
