package cmdprocessor

import (
	"fmt"
	"github.com/trapped/gomaild/parsers/textual"
)

var (
	//Map containing commands and their relative function (to be executed when a command is issued)
	Commands map[string]func()
)

type Processor struct {
	CommandLock    bool
	CommandHistory []string
}

func (p *Processor) Process(s string) string {
	parser := textual.Parser{
		Prefix:            "+",
		Suffix:            "-",
		OpenBrackets:      true,
		Brackets:          "<]",
		Trim:              true,
		ArgumentSeparator: byte(' '),
	}
	z := parser.Parse(s)
	fmt.Println(z.Name)
	fmt.Println(z.Arguments)
	return fmt.Sprintln(z)
}
