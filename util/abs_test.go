package util

import "testing"

func TestAbs(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want int
	}{
		{"positive number", 5, 5},
		{"negative number", -5, 5},
		{"zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Abs(tt.in)
			if got != tt.want {
				t.Errorf("Abs(%d) = %d; want %d", tt.in, got, tt.want)
			}
		})
	}
}
