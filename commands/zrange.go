package commands

import (
	"fmt"
	"math"
	"strings"
)

// ZRANGE key start stop [BYSCORE | BYLEX] [REV] [WITHSCORES]

type ZRangeCmd struct {
	Key string

	StartIndex int
	StopIndex  int

	MinScore float64
	MaxScore float64

	MinKey string
	MaxKey string

	ByIndex bool
	ByScore bool
	ByLex   bool

	Reverse bool

	WithScores bool
}

func (cmd *ZRangeCmd) Build(args CmdArgs) error {
	if len(args) < 3 {
		return errWrongNumberOfArgs
	}

	cmd.Key = args[0]

	if len(args) > 3 {
		for _, arg := range args[2:] {
			switch strings.ToLower(arg) {
			case "byscore":
				cmd.ByScore = true
			case "bylex":
				cmd.ByLex = true
			case "rev":
				cmd.Reverse = true
			case "withscores":
				cmd.WithScores = true
			}
		}
	}

	if cmd.ByScore {
		minscore, maxscore := cmd.parseFloatScoreRange(args[1], args[2])
		cmd.MinScore, cmd.MaxScore = minscore, maxscore
		return nil
	}

	if cmd.ByLex {
		cmd.MinKey = args[1]
		cmd.MaxKey = args[2]
		return nil
	}

	cmd.ByIndex = true
	fmt.Sscanf(args[1], "%d", &cmd.StartIndex)
	fmt.Sscanf(args[2], "%d", &cmd.StopIndex)

	return nil
}

func (cmd *ZRangeCmd) parseFloatScoreRange(arg1, arg2 string) (minScore float64, maxScore float64) {
	minScore = -math.MaxFloat64
	maxScore = math.MaxFloat64
	if arg1 != "-inf" && arg1 != "+inf" {
		if arg1[0] == '(' {
			fmt.Sscanf(arg1, "(%f", &minScore)
			minScore = math.Nextafter(minScore, math.Inf(1))
		} else {
			fmt.Sscanf(arg1, "%f", &minScore)
		}
	}

	if arg2 != "-inf" && arg2 != "+inf" {
		if arg2[0] == '(' {
			fmt.Sscanf(arg2, "(%f", &maxScore)
			maxScore = math.Nextafter(maxScore, math.Inf(-1))
		} else {
			fmt.Sscanf(arg2, "%f", &maxScore)
		}
	}

	return minScore, maxScore
}
