package main

import "testing"

func TestCompass_Rotate(t *testing.T) {
	tests := []struct {
		name   string
		c      Compass
		rotate int
		want   Compass
	}{
		{
			name:   "Stay still",
			c:      N,
			rotate: 0,
			want:   N,
		}, {
			name:   "Single step",
			c:      N,
			rotate: 1,
			want:   NE,
		},{
			name:   "Around the compass",
			c:      NW,
			rotate: 1,
			want:   N,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rotate := tt.c.Rotate(tt.rotate)
			if got := rotate; got != tt.want {
				t.Errorf("Rotate() = %v, want %v", got, tt.want)
			}
		})
	}
}
