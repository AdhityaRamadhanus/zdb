package commands

import (
	"strconv"
	"strings"
)

// ZADD key [NX | XX] [GT | LT] [CH] [INCR] score member [score member...]

// ZADD supports a list of options, specified after the name of the key and before the first score argument. Options are:

// XX: Only update elements that already exist. Don't add new elements.
// NX: Only add new elements. Don't update already existing elements.
// LT: Only update existing elements if the new score is less than the current score. This flag doesn't prevent adding new elements.
// GT: Only update existing elements if the new score is greater than the current score. This flag doesn't prevent adding new elements.
// CH: Modify the return value from the number of new elements added, to the total number of elements changed (CH is an abbreviation of changed). Changed elements are new elements added and elements already existing for which the score was updated. So elements specified in the command line having the same score as they had in the past are not counted. Note: normally the return value of ZADD only counts the number of new elements added.
// INCR: When this option is specified ZADD acts like ZINCRBY. Only one score-element pair can be specified in this mode.
// Note: The GT, LT and NX options are mutually exclusive.

type ZMember struct {
	Score float64
	Key   string
}
type ZADDCmd struct {
	Key string
	XX  bool
	NX  bool

	LT bool
	GT bool

	CH   bool
	INCR bool

	Members []ZMember
}

func (z *ZADDCmd) Build(args CmdArgs) error {
	if len(args) < 3 {
		return errWrongNumberOfArgs
	}

	z.Key = args[0]
	cutOff := 0
	for i := 1; i < len(args); {
		isArgs := true
		switch strings.ToLower(args[i]) {
		case "xx":
			z.XX = true
		case "nx":
			z.NX = true
		case "lt":
			z.LT = true
			z.NX = false
		case "gt":
			z.GT = true
			z.NX = false
		case "ch":
			z.CH = true
		case "incr":
			z.INCR = true
		default:
			isArgs = false
		}

		if !isArgs {
			cutOff = i
			break
		}
		i++
	}

	args = args[cutOff:]

	if len(args)%2 != 0 {
		return errWrongNumberOfArgs
	}

	// get Z members
	for i := 0; i < len(args)/2; i++ {
		score, err := strconv.ParseFloat(args[i*2], 64)
		if err != nil {
			continue
		}
		member := args[i*2+1]
		z.Members = append(z.Members, ZMember{
			Score: score,
			Key:   member,
		})
	}

	if len(z.Members) == 0 {
		return errWrongNumberOfArgs
	}

	return nil
}
