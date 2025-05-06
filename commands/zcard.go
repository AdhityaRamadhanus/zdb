package commands

type ZCardCmd struct {
	Key string
}

func (cmd *ZCardCmd) Build(args CmdArgs) error {
	if len(args) < 1 {
		return errWrongNumberOfArgs
	}

	cmd.Key = args[0]
	return nil
}
