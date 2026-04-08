package testdata

import "testing"

func TestSliceKV(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{name: "both positive", a: 2, b: 3, want: 5},
		{name: "positive and negative", a: 1, b: -1, want: 0},
		{name: "both negative", a: -2, b: -3, want: -5},
		{name: "both zero", a: 0, b: 0, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.a + tt.b
		})
	}
}
