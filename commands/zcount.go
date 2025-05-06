package commands

import "strconv"

type ZCountCmd struct {
	Key string
	Min float64
	Max float64
}

func (cmd *ZCountCmd) Build(args CmdArgs) (err error) {
	if len(args) < 3 {
		return errWrongNumberOfArgs
	}

	cmd.Key = args[0]
	cmd.Min, err = strconv.ParseFloat(args[1], 64)
	if err != nil {
		return err
	}
	cmd.Max, err = strconv.ParseFloat(args[2], 64)
	if err != nil {
		return err
	}

	return nil
}
