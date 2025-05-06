package commands

import (
	"strconv"
	"strings"
)

// SCAN cursor [MATCH pattern] [COUNT count]

type ScanCmd struct {
	Cursor       string
	MatchPattern string
	Count        int
}

func (cmd *ScanCmd) Build(args CmdArgs) error {
	if len(args) < 1 {
		return errWrongNumberOfArgs
	}

	cmd.Count = 10
	cmd.MatchPattern = ""
	cmd.Cursor = args[0]

	args = args[1:]

	for i := 0; i < len(args); {
		switch strings.ToLower(args[i]) {
		case "match":
			if (i + 1) >= len(args) {
				return errWrongNumberOfArgs
			}
			i++
			cmd.MatchPattern = args[i]
		case "count":
			if (i + 1) >= len(args) {
				return errWrongNumberOfArgs
			}
			i++
			cmd.Count, _ = strconv.Atoi(args[i])
		}
		i++
	}

	return nil
}
