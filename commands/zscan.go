package commands

type ZScanCmd struct {
	Key     string
	ScanCmd ScanCmd
}

func (cmd *ZScanCmd) Build(args CmdArgs) error {
	if len(args) < 2 {
		return errWrongNumberOfArgs
	}
	cmd.Key = args[0]
	return cmd.ScanCmd.Build(args[1:])
}
