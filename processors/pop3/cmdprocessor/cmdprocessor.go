package cmdprocessor

import (
	"github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/pop3/session"
)

var (
	//Map containing commands and their relative function (to be executed when a command is issued)
	Commands map[string]func(*Session, textual.Command) (string, error)
)

type Processor struct {
	Session *Session
}

func (p *Processor) Process(s string) string {
	parser := textual.Parser{
		Prefix:            "",
		Suffix:            "",
		OpenBrackets:      false,
		Brackets:          "",
		Trim:              true,
		ArgumentSeparator: byte(' '),
	}
	z := parser.Parse(s)
	if _, exists := Commands[z.Name]; exists {
		res, err := Commands[z.Name](p.Session, z)
		if err != nil {
			return "-ERR " + err.Error() + "\r\n"
		}
		return "+OK " + res + "\r\n"
	}
	return "-ERR command not found" + "\r\n"
}
