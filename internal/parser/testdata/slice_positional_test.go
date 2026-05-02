package testdata

import "testing"

func TestSlicePositional(t *testing.T) {
	tests := []struct {
		name string
		a    int
	}{
		{"first case", 1},
		{"second case", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.a
		})
	}
}
