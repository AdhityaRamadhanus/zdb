package zdb

import (
	"testing"
)

func TestMask64(t *testing.T) {
	tests := []struct {
		Name string
		num  uint
		want uint64
	}{
		{
			Name: "0",
			num:  0,
			want: 1,
		},
		{
			Name: "8",
			num:  8,
			want: 16,
		},
		{
			Name: "7",
			num:  7,
			want: 8,
		},
		{
			Name: "16",
			num:  16,
			want: 32,
		},
		{
			Name: "15",
			num:  15,
			want: 16,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			count := Mask64(test.num)
			if count != test.want {
				t.Errorf("got %v, want %v\n", count, test.want)
			}
		})
	}
}
