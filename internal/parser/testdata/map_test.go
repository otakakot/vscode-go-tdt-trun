package testdata

import "testing"

func TestMap(t *testing.T) {
	tests := map[string]struct {
		input  string
		result string
	}{
		"empty string": {
			input:  "",
			result: "",
		},
		"single char": {
			input:  "x",
			result: "x",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_ = tt.input
		})
	}
}
