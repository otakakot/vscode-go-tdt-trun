package calc_test

import (
	"testing"

	"github.com/otakakot/vscode-go-tdt-trun/example/calc"
)

func TestAdd(t *testing.T) {
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
			got := calc.Add(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSub(t *testing.T) {
	t.Run("both positive", func(t *testing.T) {
		if got := calc.Sub(2, 3); got != -1 {
			t.Errorf("Sub() = %v, want %v", got, -1)
		}
	})
	t.Run("positive and negative", func(t *testing.T) {
		if got := calc.Sub(1, -1); got != 2 {
			t.Errorf("Sub() = %v, want %v", got, 2)
		}
	})
	t.Run("both negative", func(t *testing.T) {
		if got := calc.Sub(-2, -3); got != 1 {
			t.Errorf("Sub() = %v, want %v", got, 1)
		}
	})
	t.Run("both zero", func(t *testing.T) {
		if got := calc.Sub(0, 0); got != 0 {
			t.Errorf("Sub() = %v, want %v", got, 0)
		}
	})
}
