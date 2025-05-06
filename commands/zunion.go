package commands

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var (
	errZUnionNeed2Keys        = errors.New("zunion command need at least 2 keys")
	errKeysDoesntMatchNumKeys = errors.New("number of keys doesnt match numkeys")
	errWeightsDoesntMatchKeys = errors.New("number of weights doesnt match number of keys")
	errUnsupportedAggFunc     = errors.New("unsupported aggregate function")
)

type ZUnionCmd struct {
	NumKeys    int
	Keys       []string
	Weights    []float64
	Aggregate  string
	WithScores bool
}

// ZUNION numkeys key [key ...] [WEIGHTS weight [weight ...]] [AGGREGATE <SUM | MIN | MAX>] [WITHSCORES]

func (cmd *ZUnionCmd) Build(args CmdArgs) (err error) {
	// at least 3 args, numkeys and 2 keys
	if len(args) < 3 {
		return errWrongNumberOfArgs
	}

	cmd.NumKeys, err = strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	if cmd.NumKeys < 2 {
		return errZUnionNeed2Keys
	}

	if (1 + cmd.NumKeys) > len(args) {
		return errKeysDoesntMatchNumKeys
	}

	cmd.Keys = append(cmd.Keys, args[1:(1+cmd.NumKeys)]...)
	args = args[(1 + cmd.NumKeys):]
	for i := 0; i < len(args); {
		switch strings.ToLower(args[i]) {
		case "weights":
			if i+cmd.NumKeys >= len(args) {
				return errWeightsDoesntMatchKeys
			}
			for i := i + 1; i <= cmd.NumKeys; i++ {
				weight, err := strconv.ParseFloat(args[i], 64)
				if err != nil {
					return errWrongNumberOfArgs
				}
				cmd.Weights = append(cmd.Weights, weight)
			}
		case "aggregate":
			i++
			if i >= len(args) {
				return errWeightsDoesntMatchKeys
			}

			aggFunc := strings.ToLower(args[i])
			fmt.Println("Agg func is ", aggFunc)
			if slices.Index([]string{"sum", "min", "max"}, aggFunc) == -1 {
				return errUnsupportedAggFunc
			}

			cmd.Aggregate = aggFunc
		case "withscores":
			cmd.WithScores = true
		}

		i++
	}

	if len(cmd.Weights) == 0 {
		cmd.Weights = slices.Repeat([]float64{1}, cmd.NumKeys)
	}

	return nil
}
