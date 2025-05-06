package commands

import "errors"

var (
	errWrongNumberOfArgs = errors.New("wrong number of arguments")
)

type CmdAble interface {
	Build(CmdArgs) error
}

type CmdArgs []string
