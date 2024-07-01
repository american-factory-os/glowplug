package sparkplug

import "testing"

func TestNextSequenceNumber(t *testing.T) {
	tests := []struct {
		name    string
		current uint64
		want    uint64
	}{
		{
			name:    "Test sequence number 0",
			current: 0,
			want:    1,
		},
		{
			name:    "Test sequence number 255",
			current: 255,
			want:    0,
		},
		{
			name:    "Test sequence number 254",
			current: 254,
			want:    255,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NextSequenceNumber(tt.current); got != tt.want {
				t.Errorf("NextSequenceNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
