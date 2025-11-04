package params

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBool(t *testing.T) {
	type want struct {
		Value   bool
		Present bool
	}

	type Test struct {
		Field want `json:"field"`
		Value want `json:"value"`
	}

	type result struct {
		Field Bool `json:"field"`
		Value Bool `json:"value"`
	}

	tests := []struct {
		name    string
		input   string
		output  string
		want    Test
		wantErr bool
	}{
		{
			name:  "Valid JSON with boolean",
			input: `{"field":true,"value":false}`,
			want: Test{
				Field: want{Value: true, Present: true},
				Value: want{Value: false, Present: true},
			},
			wantErr: false,
		},
		{
			name:   "Valid JSON with quoted boolean",
			input:  `{"field":"true","value":"true"}`,
			output: `{"field":true,"value":true}`,
			want: Test{
				Field: want{Value: true, Present: true},
				Value: want{Value: true, Present: true},
			},
			wantErr: false,
		},
		{
			name:    "Empty JSON",
			input:   `{}`,
			output:  `{"field":null,"value":null}`,
			want:    Test{},
			wantErr: false,
		},
		{
			name:   "Null JSON",
			input:  `{"field":null,"value":null}`,
			output: `{"field":null,"value":null}`,
			want: Test{
				Field: want{Present: false},
				Value: want{Present: false},
			},
			wantErr: false,
		},
		{
			name:    "Invalid JSON",
			input:   `{"field": 123,"value": 456`,
			want:    Test{},
			wantErr: true,
		},
		{
			name:    "Invalid JSON with extra comma",
			input:   `{"field": true,"value": false,	}`,
			want:    Test{},
			wantErr: true,
		},
		{
			name:   "Missing field",
			input:  `{"value":false}`,
			output: `{"field":null,"value":false}`,
			want: Test{
				Field: want{Present: false},
				Value: want{Value: false, Present: true},
			},
			wantErr: false,
		},
		{
			name:   "Missing value",
			input:  `{"field":true}`,
			output: `{"field":true,"value":null}`,
			want: Test{
				Field: want{Value: true, Present: true},
				Value: want{Present: false},
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
				require.NoError(t, err, "Marshal should not return an error")
				require.JSONEq(t, tt.output, string(js), "Marshalled JSON should match the original input")
			}
		})
	}
}
