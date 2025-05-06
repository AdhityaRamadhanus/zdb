package commands

type ZUnionStoreCmd struct {
	DstKey    string
	ZUnionCmd ZUnionCmd
}

func (cmd *ZUnionStoreCmd) Build(args CmdArgs) error {
	if len(args) < 4 {
		return errWrongNumberOfArgs
	}

	cmd.DstKey = args[0]
	args = args[1:]

	cmd.ZUnionCmd.WithScores = true
	return cmd.ZUnionCmd.Build(args)
}
