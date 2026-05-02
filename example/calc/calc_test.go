package calc_test

import (
	"testing"

	"github.com/otakakot/vscode-go-tdt-trun/example/calc"
)

func TestAdd_SliceKeyValue(t *testing.T) {
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

func TestAdd_SlicePositional(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{"both positive", 2, 3, 5},
		{"positive and negative", 1, -1, 0},
		{"both negative", -2, -3, -5},
		{"both zero", 0, 0, 0},
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

func TestAdd_Map(t *testing.T) {
	tests := map[string]struct {
		a    int
		b    int
		want int
	}{
		"both positive":         {a: 2, b: 3, want: 5},
		"positive and negative": {a: 1, b: -1, want: 0},
		"both negative":         {a: -2, b: -3, want: -5},
		"both zero":             {a: 0, b: 0, want: 0},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := calc.Add(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdd_Literal(t *testing.T) {
	t.Run("both positive", func(t *testing.T) {
		if got := calc.Add(2, 3); got != 5 {
			t.Errorf("Add() = %v, want %v", got, 5)
		}
	})
	t.Run("both zero", func(t *testing.T) {
		if got := calc.Add(0, 0); got != 0 {
			t.Errorf("Add() = %v, want %v", got, 0)
		}
	})
}
