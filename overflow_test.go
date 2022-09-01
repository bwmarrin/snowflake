package snowflake

import (
	"testing"
	"time"
)

func TestOverflow_IsZero(t *testing.T) {
	tests := []struct {
		name string
		o    Overflow
		want bool
	}{
		{
			name: "zero",
			o:    Overflow{},
			want: true,
		},
		{
			name: "non-zero",
			o: Overflow{
				Time: time.Now(),
				Ms:   1,
				Step: 2,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOverflow_Duration(t *testing.T) {
	tests := []struct {
		name string
		o    Overflow
		want time.Duration
	}{
		{
			name: "zero",
			o:    Overflow{},
			want: 0,
		},
		{
			name: "non-zero",
			o: Overflow{
				Time: time.Now(),
				Ms:   1,
				Step: 2,
			},
			want: 1 * time.Millisecond,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.Duration(); got != tt.want {
				t.Errorf("Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOverflow_UntilCleared(t *testing.T) {
	tests := []struct {
		name string
		o    Overflow
	}{
		{
			name: "zero",
			o:    Overflow{},
		},
		{
			name: "non-zero",
			o: Overflow{
				Time: time.Now(),
				Ms:   1,
				Step: 5,
			},
		},
		{
			name: "non-zero, old",
			o: Overflow{
				Time: time.Now().Add(-100 * time.Millisecond),
				Ms:   1,
				Step: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.o.Duration()
			dc := tt.o.DurationToClear()
			switch {
			case d == 0:
				if dc != 0 {
					t.Fatal("DurationToClear must be zero when Duration is zero")
				}
			default:
				if dc >= d {
					t.Fatal("DurationToClear must be less than Duration")
				}
			}
		})
	}
}
