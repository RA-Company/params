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
	Value   time.Time `json:"-"`
	Present bool      `json:"-"` // Present indicates if the time is present or not
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
	dst.Value = time.Time{}
	if len(data) == 0 || string(data) == "null" {
		dst.Present = false
		return nil
	}

	dst.Present = true

	for _, layout := range timeLayouts {
		if t, err := time.Parse(layout, strings.Trim(string(data), `"`)); err == nil {
			dst.Value = t
			return nil
		}
	}

	return fmt.Errorf("invalid time format: %s", string(data))
}

// MarshalJSON implements the json.Marshaler interface.
// It returns "null" if the time is not present.
//
// Returns:
//   - []byte: JSON representation of the time.
//   - error: An error if marshaling fails, otherwise nil.
func (dst *Time) MarshalJSON() ([]byte, error) {
	if !dst.Present {
		return []byte("null"), nil
	}
	return dst.Value.MarshalJSON()
}

// IsZero checks if the Time is zero or not present.
//
// Returns:
//   - bool: true if the Time is zero or not present, otherwise false.
func (dst *Time) IsZero() bool {
	return !dst.Present || dst.Value.IsZero()
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
	if !dst.Present {
		return ""
	}
	return dst.Value.Format(layout)
}

// Equal checks if two Time instances are equal.
//
// Returns:
//   - bool: true if both Time instances are equal, otherwise false.
func (dst *Time) String() string {
	if !dst.Present {
		return "null"
	}
	return dst.Value.String()
}
