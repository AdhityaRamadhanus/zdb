//go:build unit

package commands

import (
	"slices"
	"testing"
)

func TestZUnionCmd(t *testing.T) {
	tests := []struct {
		Name    string
		Args    []string
		Want    ZUnionCmd
		WantErr error
	}{
		{
			Name: "All options used and correct",
			Args: []string{"2", "zset1", "zset2", "weights", "0.5", "0.7", "aggregate", "max", "WitHScores"},
			Want: ZUnionCmd{
				NumKeys:    2,
				Keys:       []string{"zset1", "zset2"},
				Weights:    []float64{0.5, 0.7},
				WithScores: true,
				Aggregate:  "max",
			},
			WantErr: nil,
		},
		{
			Name:    "Incorrect number of weights provided",
			Args:    []string{"2", "zset1", "zset2", "weights", "0.5"},
			Want:    ZUnionCmd{},
			WantErr: errWeightsDoesntMatchKeys,
		},
		{
			Name:    "Insufficient number of args",
			Args:    []string{"2", "zset1"},
			Want:    ZUnionCmd{},
			WantErr: errWrongNumberOfArgs,
		},
		{
			Name:    "Incorrect number of keys",
			Args:    []string{"3", "zset1", "zset2"},
			Want:    ZUnionCmd{},
			WantErr: errKeysDoesntMatchNumKeys,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got := ZUnionCmd{}
			err := got.Build(test.Args)
			if err != test.WantErr {
				t.Errorf("got err %v, want err %v", err, test.WantErr)
			} else {
				if test.WantErr == nil {
					checkZUnionCmd(t, got, test.Want)
				}
			}
		})
	}
}

func checkZUnionCmd(t *testing.T, got, want ZUnionCmd) {
	t.Helper()

	isSameKeys := slices.Compare(got.Keys, want.Keys) == 0
	isSameWeights := slices.Compare(got.Weights, want.Weights) == 0

	isAllSame := (got.NumKeys == want.NumKeys &&
		got.Aggregate == want.Aggregate &&
		isSameKeys &&
		isSameWeights &&
		got.WithScores == want.WithScores)

	if !isAllSame {
		t.Errorf("got %v, want %v\n", got, want)
	}
}
