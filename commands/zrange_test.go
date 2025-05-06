//go:build unit

package commands

import (
	"math"
	"testing"
)

// ZRANGE key start stop [BYSCORE | BYLEX] [REV] [WITHSCORES]

func TestZRange(t *testing.T) {
	tests := []struct {
		Name    string
		Args    []string
		Want    ZRangeCmd
		WantErr error
	}{
		{
			Name: "zrange by index with scores",
			Args: []string{"zset1", "0", "-1", "WitHScores"},
			Want: ZRangeCmd{
				Key:        "zset1",
				StartIndex: 0,
				StopIndex:  -1,
				ByLex:      false,
				ByScore:    false,
				ByIndex:    true,
				WithScores: true,
			},
			WantErr: nil,
		},
		{
			Name: "zrange by score without scores",
			Args: []string{"zset1", "50", "70.5", "byscore"},
			Want: ZRangeCmd{
				Key:        "zset1",
				MinScore:   50,
				MaxScore:   70.5,
				ByLex:      false,
				ByScore:    true,
				ByIndex:    false,
				WithScores: false,
			},
			WantErr: nil,
		},
		{
			Name: "zrange by score with exclusive range without scores",
			Args: []string{"zset1", "(50", "(70.5", "byscore"},
			Want: ZRangeCmd{
				Key:        "zset1",
				MinScore:   math.Nextafter(50, math.Inf(1)),
				MaxScore:   math.Nextafter(70.5, math.Inf(-1)),
				ByLex:      false,
				ByScore:    true,
				ByIndex:    false,
				WithScores: false,
			},
			WantErr: nil,
		},
		{
			Name: "zrange by lex with scores",
			Args: []string{"zset1", "A", "Z", "bylex", "withscores"},
			Want: ZRangeCmd{
				Key:        "zset1",
				MinKey:     "A",
				MaxKey:     "Z",
				ByLex:      true,
				ByScore:    false,
				ByIndex:    false,
				WithScores: true,
			},
			WantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got := ZRangeCmd{}
			err := got.Build(test.Args)
			if err != test.WantErr {
				t.Errorf("got err %v, want err %v", err, test.WantErr)
			} else {
				if test.WantErr == nil {
					checkZRangeCmd(t, got, test.Want)
				}
			}
		})
	}
}

func checkZRangeCmd(t *testing.T, got, want ZRangeCmd) {
	t.Helper()

	if !(got == want) {
		t.Errorf("got %+v, want %+v\n", got, want)
	}
}
