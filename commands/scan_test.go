//go:build unit

package commands

import (
	"testing"
)

func TestScan(t *testing.T) {
	tests := []struct {
		Name    string
		Args    []string
		Want    ScanCmd
		WantErr error
	}{
		{
			Name: "All options used and correct",
			Args: []string{"0", "match", "db:*", "count", "10"},
			Want: ScanCmd{
				MatchPattern: "db:*",
				Cursor:       "0",
				Count:        10,
			},
			WantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got := ScanCmd{}
			err := got.Build(test.Args)
			if err != test.WantErr {
				t.Errorf("got err %v, want err %v", err, test.WantErr)
			} else {
				if test.WantErr == nil {
					if got != test.Want {
						t.Errorf("got %v, want %v", got, test.Want)
					}
				}
			}
		})
	}
}
