package commands

type ZInterStoreCmd struct {
	DstKey    string
	ZInterCmd ZInterCmd
}

func (cmd *ZInterStoreCmd) Build(args CmdArgs) error {
	if len(args) < 4 {
		return errWrongNumberOfArgs
	}

	cmd.DstKey = args[0]
	args = args[1:]

	cmd.ZInterCmd.WithScores = true
	return cmd.ZInterCmd.Build(args)
}
