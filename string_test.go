package params

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	type Test struct {
		Field String `json:"field"`
		Value String `json:"value"`
	}

	tests := []struct {
		name    string
		input   string
		output  string
		want    Test
		wantErr bool
	}{
		{
			name:  "Valid JSON with quotes",
			input: `{"field":"test\"Field","value":"test\"Value"}`,
			want: Test{
				Field: String{Value: `test"Field`, Present: true},
				Value: String{Value: `test"Value`, Present: true},
			},
			wantErr: false,
		},
		{
			name:  "Valid JSON without quotes",
			input: `{"field":"testField","value":"testValue"}`,
			want: Test{
				Field: String{Value: "testField", Present: true},
				Value: String{Value: "testValue", Present: true},
			},
			wantErr: false,
		},
		{
			name:    "Empty JSON",
			input:   `{}`,
			output:  `{"field":"","value":""}`,
			want:    Test{},
			wantErr: false,
		},
		{
			name:    "Null JSON",
			input:   `{"field":null,"value":null}`,
			output:  `{"field":"","value":""}`,
			want:    Test{Field: String{Present: false}, Value: String{Present: false}},
			wantErr: false,
		},
		{
			name:    "Invalid JSON",
			input:   `{"field": "testField","value": "testValue"`,
			want:    Test{},
			wantErr: true,
		},
		{
			name:    "Invalid JSON with extra comma",
			input:   `{"field": "testField","value": "testValue",	}`,
			want:    Test{},
			wantErr: true,
		},
		{
			name:   "Missing field",
			input:  `{"value":"testValue"}`,
			output: `{"field":"","value":"testValue"}`,
			want: Test{
				Field: String{Present: false},
				Value: String{Value: "testValue", Present: true},
			},
			wantErr: false,
		},
		{
			name:   "Missing value",
			input:  `{"field":"testField"}`,
			output: `{"field":"testField","value":""}`,
			want: Test{
				Field: String{Value: "testField", Present: true},
				Value: String{Present: false},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.output == "" {
				tt.output = tt.input
			}
			var test Test
			err := json.Unmarshal([]byte(tt.input), &test)
			if tt.wantErr {
				require.Error(t, err, "Unmarshal should return an error")
			} else {
				require.NoError(t, err, "Unmarshal should not return an error")
				require.Equal(t, tt.want.Field.Value, test.Field.Value, "Field value should match the input")
				require.Equal(t, tt.want.Field.Present, test.Field.Present, "Field should be present")
				require.Equal(t, tt.want.Value.Value, test.Value.Value, "Value should match the input")
				require.Equal(t, tt.want.Value.Present, test.Value.Present, "Value should be present")

				js, err := json.Marshal(test)
				require.NoError(t, err, "Marshal should not return an error")
				require.JSONEq(t, tt.output, string(js), "Marshalled JSON should match the original input")
			}
		})
	}
}
