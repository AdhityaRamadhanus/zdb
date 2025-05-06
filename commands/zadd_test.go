//go:build unit

package commands

import (
	"slices"
	"testing"
)

func TestZAdd(t *testing.T) {
	tests := []struct {
		Name    string
		Args    []string
		Want    ZADDCmd
		WantErr error
	}{
		{
			Name: "All options used and correct",
			Args: []string{"zset1", "xx", "LT", "INCR", "85", "MemberA"},
			Want: ZADDCmd{
				Key:  "zset1",
				NX:   false,
				XX:   true,
				LT:   true,
				GT:   false,
				INCR: true,
				Members: []ZMember{
					{
						Score: 85,
						Key:   "MemberA",
					},
				},
			},
			WantErr: nil,
		},
		{
			Name: "Insufficient arg for members",
			Args: []string{"zset1", "85", "MemberA", "78"},
			Want: ZADDCmd{
				Key:     "zset1",
				NX:      false,
				XX:      false,
				LT:      false,
				GT:      false,
				INCR:    false,
				Members: []ZMember{},
			},
			WantErr: errWrongNumberOfArgs,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got := ZADDCmd{}
			err := got.Build(test.Args)
			if err != test.WantErr {
				t.Errorf("got err %v, want err %v", err, test.WantErr)
			} else {
				if test.WantErr == nil {
					checkZAddCmd(t, got, test.Want)
				}
			}
		})
	}
}

func checkZAddCmd(t *testing.T, got, want ZADDCmd) {
	t.Helper()

	isSameMembers := slices.CompareFunc(got.Members, want.Members, func(member1, member2 ZMember) int {
		if member1.Score < member2.Score || (member1.Score == member2.Score && member1.Key < member2.Key) {
			return -1
		} else if member1.Score > member2.Score || (member1.Score == member2.Score && member1.Key > member2.Key) {
			return 1
		} else {
			return 0
		}
	}) == 0

	isAllSame := (got.Key == want.Key &&
		got.XX == want.XX &&
		got.NX == want.NX &&
		got.LT == want.LT &&
		got.GT == want.GT &&
		got.CH == want.CH &&
		got.INCR == want.INCR &&
		isSameMembers)

	if !isAllSame {
		t.Errorf("got %+v, want %+v\n", got, want)
	}
}
