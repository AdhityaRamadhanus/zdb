package commands

type ZRemCmd struct {
	Key     string
	Members []string
}

// ZREM key member [member ...]
// RESP2/RESP3 Reply
// Integer reply: the number of members removed from the sorted set, not including non-existing members.

func (cmd *ZRemCmd) Build(args CmdArgs) error {
	if len(args) < 2 {
		return errWrongNumberOfArgs
	}

	cmd.Key = args[0]
	cmd.Members = append(cmd.Members, args[1:]...)
	return nil
}
