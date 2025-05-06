package commands

import (
	"errors"
	"strconv"
	"strings"
)

var (
	errZDiffNeed2Keys = errors.New("zdiff command need at least 2 keys")
)

type ZDiffCmd struct {
	NumKeys    int
	Keys       []string
	WithScores bool
}

// ZDIFF numkeys key [key ...] [WITHSCORES]

func (cmd *ZDiffCmd) Build(args CmdArgs) (err error) {
	// at least 3 args, numkeys and 2 keys
	if len(args) < 3 {
		return errWrongNumberOfArgs
	}

	cmd.NumKeys, err = strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	if cmd.NumKeys < 2 {
		return errZDiffNeed2Keys
	}

	if (1 + cmd.NumKeys) > len(args) {
		return errKeysDoesntMatchNumKeys
	}

	cmd.Keys = append(cmd.Keys, args[1:(1+cmd.NumKeys)]...)
	args = args[(1 + cmd.NumKeys):]
	for i := 0; i < len(args); {
		switch strings.ToLower(args[i]) {
		case "withscores":
			cmd.WithScores = true
		}

		i++
	}
	return nil
}
