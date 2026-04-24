package params

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUUID(t *testing.T) {
	type want struct {
		Value   uuid.UUID
		Present bool
	}
	type Test struct {
		Field want `json:"field"`
		Value want `json:"value"`
	}
	type result struct {
		Field UUID `json:"field"`
		Value UUID `json:"value"`
	}

	tests := []struct {
		name    string
		input   string
		output  string
		want    Test
		wantErr bool
	}{
		{
			name:  "Valid JSON",
			input: `{"field":"71d38b17-3610-4e7c-8e9c-11725aa5b133","value":"7332765a-c5d0-4dd3-a62a-4f3df8c21c25"}`,
			want: Test{
				Field: want{Value: uuid.MustParse("71d38b17-3610-4e7c-8e9c-11725aa5b133"), Present: true},
				Value: want{Value: uuid.MustParse("7332765a-c5d0-4dd3-a62a-4f3df8c21c25"), Present: true},
			},
			wantErr: false,
		},
		{
			name:   "Empty JSON",
			input:  `{}`,
			output: `{"field":"","value":""}`,
			want: Test{
				Field: want{Present: false},
				Value: want{Present: false},
			},
			wantErr: false,
		},
		{
			name:    "Null JSON",
			input:   `{"field":null,"value":null}`,
			output:  `{"field":"","value":""}`,
			want:    Test{Field: want{Present: false}, Value: want{Present: false}},
			wantErr: false,
		},
		{
			name:    "Invalid JSON",
			input:   `{"field": "testField","value": "testValue"`,
			want:    Test{},
			wantErr: true,
		},
		{
			name:   "Missing field",
			input:  `{"value":"7332765a-c5d0-4dd3-a62a-4f3df8c21c25"}`,
			output: `{"field":"","value":"7332765a-c5d0-4dd3-a62a-4f3df8c21c25"}`,
			want: Test{
				Field: want{Present: false},
				Value: want{Value: uuid.MustParse("7332765a-c5d0-4dd3-a62a-4f3df8c21c25"), Present: true},
			},
			wantErr: false,
		},
		{
			name:   "Missing value",
			input:  `{"field":"71d38b17-3610-4e7c-8e9c-11725aa5b133"}`,
			output: `{"field":"71d38b17-3610-4e7c-8e9c-11725aa5b133","value":""}`,
			want: Test{
				Field: want{Value: uuid.MustParse("71d38b17-3610-4e7c-8e9c-11725aa5b133"), Present: true},
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
