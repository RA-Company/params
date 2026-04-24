package params

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFloat(t *testing.T) {
	type want struct {
		Value   float64
		Present bool
	}

	type Test struct {
		Field want `json:"field"`
		Value want `json:"value"`
	}

	type result struct {
		Field Float `json:"field"`
		Value Float `json:"value"`
	}

	tests := []struct {
		name    string
		input   string
		output  string
		want    Test
		wantErr bool
	}{
		{
			name:  "Valid JSON with float",
			input: `{"field":123.45,"value":456.78}`,
			want: Test{
				Field: want{Value: 123.45, Present: true},
				Value: want{Value: 456.78, Present: true},
			},
			wantErr: false,
		},
		{
			name:   "Valid JSON with quoted float",
			input:  `{"field":"123.45","value":"456.78"}`,
			output: `{"field":123.45,"value":456.78}`,
			want: Test{
				Field: want{Value: 123.45, Present: true},
				Value: want{Value: 456.78, Present: true},
			},
			wantErr: false,
		},
		{
			name:    "Empty JSON",
			input:   `{}`,
			output:  `{"field":0,"value":0}`,
			want:    Test{},
			wantErr: false,
		},
		{
			name:    "Null JSON",
			input:   `{"field":null,"value":null}`,
			output:  `{"field":0,"value":0}`,
			want:    Test{Field: want{Present: false}, Value: want{Present: false}},
			wantErr: false,
		},
		{
			name:    "Invalid JSON",
			input:   `{"field": 123,"value": 456`,
			want:    Test{},
			wantErr: true,
		},
		{
			name:    "Invalid integer value",
			input:   `{"field":"abc","value":"def"}`,
			want:    Test{},
			wantErr: true,
		},
		{
			name:    "Invalid JSON with extra comma",
			input:   `{"field": 123,"value": 456,	}`,
			want:    Test{},
			wantErr: true,
		},
		{
			name:   "Missing field",
			input:  `{"value":456}`,
			output: `{"field":0,"value":456}`,
			want: Test{
				Field: want{Present: false},
				Value: want{Value: 456, Present: true},
			},
			wantErr: false,
		},
		{
			name:   "Missing value",
			input:  `{"field":123}`,
			output: `{"field":123,"value":0}`,
			want: Test{
				Field: want{Value: 123, Present: true},
				Value: want{Present: false},
			},
			wantErr: false,
		},
		{
			name:  "Min and max float values",
			input: `{"field":-1.7976931348623157e+308,"value":1.7976931348623157e+308}`,
			want: Test{
				Field: want{Value: -1.7976931348623157e+308, Present: true},
				Value: want{Value: 1.7976931348623157e+308, Present: true},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.output == "" {
				tt.output = tt.input
			}
			var test result
			err := json.Unmarshal([]byte(tt.input), &test)
			if tt.wantErr {
				require.Error(t, err, "Unmarshal should return an error")
			} else {
				require.NoError(t, err, "Unmarshal should not return an error")
				require.Equal(t, tt.want.Field.Value, test.Field.Value(), "Field value should match the input")
				require.Equal(t, tt.want.Field.Present, test.Field.Present(), "Field should be present")
				require.Equal(t, tt.want.Value.Value, test.Value.Value(), "Value should match the input")
				require.Equal(t, tt.want.Value.Present, test.Value.Present(), "Value should be present")

				js, err := json.Marshal(test)
				fmt.Printf("Marshalled JSON: %s\n", string(js))
				require.NoError(t, err, "Marshal should not return an error")
				require.JSONEq(t, tt.output, string(js), "Marshalled JSON should match the original input")
			}
		})
	}
}
