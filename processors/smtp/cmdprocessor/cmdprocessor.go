//Package cmdprocessor provides a set of structs, variables and methods to process SMTP clients' commands.
package cmdprocessor

import (
	"github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/smtp/reply"
	. "github.com/trapped/gomaild/processors/smtp/session"
	"strings"
)

var (
	//Map containing commands and their relative function (to be executed when a command is issued)
	Commands map[string]func(*Session, textual.Statement) Reply = make(map[string]func(*Session, textual.Statement) Reply, 0)
)

//Processor is a struct to provide a throw-away command processor and session for SMTP.
type Processor struct {
	Session *Session
}

//Process processes a SMTP command and returns a result.
func (p *Processor) Process(s string) string {
	//Prepare a textual parser.
	parser := textual.Parser{
		Prefix:             "",
		Suffix:             "",
		OpenBrackets:       false,
		Brackets:           "",
		Trim:               true,
		ArgumentSeparators: []byte{' '},
	}

	//Parse the command with the parser.
	z := parser.Parse(s)

	//Run the processor for the command issued by the client (if it exists) and return the result with the correct SMTP result prefix.
	if _, exists := Commands[strings.ToLower(z.Name)]; exists {
		res := Commands[strings.ToLower(z.Name)](p.Session, z)
		if res.Code == 0 {
			return res.Message
		}
		return res.String()
	}
	return "504 command not found"
}
