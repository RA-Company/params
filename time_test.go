package params

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		data    string
		present bool
		wantErr bool
		result  time.Time
	}{
		{name: "null value", data: "null", wantErr: false, present: false, result: time.Time{}},
		{name: "valid time", data: `"2023-10-05T14:48:00Z"`, wantErr: false, present: true, result: time.Date(2023, 10, 5, 14, 48, 0, 0, time.UTC)},
		{name: "invalid time", data: `"invalid-time"`, wantErr: true, present: true, result: time.Time{}},
		{name: "empty string", data: "", wantErr: false, present: false, result: time.Time{}},
		{name: "json without quotes", data: `2023-10-05T14:48:00Z`, wantErr: false, present: true, result: time.Date(2023, 10, 5, 14, 48, 0, 0, time.UTC)},
		{name: "extra whitespace", data: `  "2023-10-05T14:48:00Z"  `, wantErr: true, present: true, result: time.Time{}},
		{name: "with milliseconds", data: `"2023-10-05T14:48:00.123Z"`, wantErr: false, present: true, result: time.Date(2023, 10, 5, 14, 48, 0, 123000000, time.UTC)},
		{name: "with timezone offset", data: `"2023-10-05T14:48:00+02:00"`, wantErr: false, present: true, result: time.Date(2023, 10, 5, 14, 48, 0, 0, time.FixedZone("UTC+2", 2*60*60))},
		{name: "space instead of T", data: `"2023-10-05 14:48:00"`, wantErr: false, present: true, result: time.Date(2023, 10, 5, 14, 48, 0, 0, time.UTC)},
		{name: "different layout with MST", data: `"2023-10-05T14:48:00 UTC"`, wantErr: false, present: true, result: time.Date(2023, 10, 5, 14, 48, 0, 0, time.UTC)},
		{name: "invalid layout", data: `"05-10-2023T14:48:00Z"`, wantErr: true, present: true, result: time.Time{}},
		{name: "with microseconds", data: `"2023-10-05T14:48:00.123456Z"`, wantErr: false, present: true, result: time.Date(2023, 10, 5, 14, 48, 0, 123456000, time.UTC)},
		{name: "only date", data: `"2023-10-05"`, wantErr: true, present: true, result: time.Time{}},
		{name: "with nanoseconds", data: `"2023-10-05T14:48:00.123456789Z"`, wantErr: false, present: true, result: time.Date(2023, 10, 5, 14, 48, 0, 123456789, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dst Time
			gotErr := dst.UnmarshalJSON([]byte(tt.data))
			if tt.wantErr {
				require.Error(t, gotErr, "expected an error but got none")
				require.Equal(t, tt.present, dst.Present(), "Present field mismatch")
			} else {
				require.NoError(t, gotErr, "unexpected error: %v", gotErr)
				require.Equal(t, tt.present, dst.Present(), "Present field mismatch")
				require.True(t, dst.Value().Equal(tt.result), "Value field mismatch: got %v, want %v", dst.Value, tt.result)
			}
		})
	}
}

func TestTime_MarshalJSON(t *testing.T) {
	type params struct {
		Value   time.Time
		Present bool
	}
	tests := []struct {
		name   string // description of this test case
		params params
		want   string
	}{
		{name: "present", params: params{Value: time.Date(2023, 10, 5, 14, 48, 0, 0, time.UTC), Present: true}, want: `"2023-10-05T14:48:00Z"`},
		{name: "zero time but present", params: params{Value: time.Time{}, Present: true}, want: `"0001-01-01T00:00:00Z"`},
		{name: "with milliseconds", params: params{Value: time.Date(2023, 10, 5, 14, 48, 0, 123000000, time.UTC), Present: true}, want: `"2023-10-05T14:48:00.123Z"`},
		{name: "with nanoseconds", params: params{Value: time.Date(2023, 10, 5, 14, 48, 0, 123456789, time.UTC), Present: true}, want: `"2023-10-05T14:48:00.123456789Z"`},
		{name: "with timezone offset", params: params{Value: time.Date(2023, 10, 5, 14, 48, 0, 0, time.FixedZone("UTC+2", 2*60*60)), Present: true}, want: `"2023-10-05T14:48:00+02:00"`},
		{name: "far future date", params: params{Value: time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC), Present: true}, want: `"3000-01-01T00:00:00Z"`},
		{name: "far past date", params: params{Value: time.Date(1000, 1, 1, 0, 0, 0, 0, time.UTC), Present: true}, want: `"1000-01-01T00:00:00Z"`},
		{name: "leap year date", params: params{Value: time.Date(2020, 2, 29, 12, 0, 0, 0, time.UTC), Present: true}, want: `"2020-02-29T12:00:00Z"`},
		{name: "end of year", params: params{Value: time.Date(2023, 12, 31, 23, 59, 59, 999999999, time.UTC), Present: true}, want: `"2023-12-31T23:59:59.999999999Z"`},
		{name: "start of year", params: params{Value: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Present: true}, want: `"2023-01-01T00:00:00Z"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dst Time
			dst.Set(tt.params.Value)
			got, gotErr := json.Marshal(&dst)
			require.NoError(t, gotErr, "unexpected error: %v", gotErr)
			require.Equal(t, tt.want, string(got), "MarshalJSON() mismatch")
		})
	}
}
