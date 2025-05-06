package commands

type ZDiffStoreCmd struct {
	DstKey   string
	ZDiffCmd ZDiffCmd
}

func (cmd *ZDiffStoreCmd) Build(args CmdArgs) error {
	if len(args) < 4 {
		return errWrongNumberOfArgs
	}

	cmd.DstKey = args[0]
	args = args[1:]

	cmd.ZDiffCmd.WithScores = true
	return cmd.ZDiffCmd.Build(args)
}
