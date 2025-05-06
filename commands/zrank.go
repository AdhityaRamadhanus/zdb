package commands

type ZRankCmd struct {
	Key    string
	Member string
}

func (cmd *ZRankCmd) Build(args CmdArgs) error {
	if len(args) < 2 {
		return errWrongNumberOfArgs
	}

	cmd.Key = args[0]
	cmd.Member = args[1]
	return nil
}
