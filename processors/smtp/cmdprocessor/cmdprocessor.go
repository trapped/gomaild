//Package cmdprocessor provides a set of structs, variables and methods to process SMTP clients' commands.
package cmdprocessor

import (
	"github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/smtp/reply"
	. "github.com/trapped/gomaild/processors/smtp/session"
	"strings"
	//SMTP commands
	"github.com/trapped/gomaild/processors/smtp/cmdprocessor/auth"
	"github.com/trapped/gomaild/processors/smtp/cmdprocessor/data"
	"github.com/trapped/gomaild/processors/smtp/cmdprocessor/ehlo"
	"github.com/trapped/gomaild/processors/smtp/cmdprocessor/helo"
	"github.com/trapped/gomaild/processors/smtp/cmdprocessor/mail"
	"github.com/trapped/gomaild/processors/smtp/cmdprocessor/noop"
	"github.com/trapped/gomaild/processors/smtp/cmdprocessor/quit"
	"github.com/trapped/gomaild/processors/smtp/cmdprocessor/rcpt"
	"github.com/trapped/gomaild/processors/smtp/cmdprocessor/rset"
	"github.com/trapped/gomaild/processors/smtp/cmdprocessor/starttls"
)

var (
	//Contains commands and their relative functions (to be executed when a command is issued)
	Commands map[string]func(*Session, textual.Statement) Reply = map[string]func(*Session, textual.Statement) Reply{
		"auth":     auth.Process,
		"data":     data.Process,
		"ehlo":     ehlo.Process,
		"helo":     helo.Process,
		"mail":     mail.Process,
		"noop":     noop.Process,
		"quit":     quit.Process,
		"rcpt":     rcpt.Process,
		"rset":     rset.Process,
		"starttls": starttls.Process,
	}
)

//Struct to provide a throw-away command processor and session for SMTP.
type Processor struct {
	Session     *Session //SMTP session, accessible by both the commands and the client handler
	LastCommand string   //The last successfully issued command
}

//Processes a SMTP command and returns a result.
func (p *Processor) Process(s string) string {
	//Prepare a textual parser.
	parser := textual.Parser{
		Prefix:             "",
		Suffix:             "",
		OpenBrackets:       true,
		Brackets:           []byte{'<', '>'},
		Trim:               true,
		ArgumentSeparators: []byte{' '},
	}

	//Parse the command with the parser.
	z := parser.Parse(s)

	if p.Session.State == COMPOSITION || p.Session.AuthState == AUTHWUSER || p.Session.AuthState == AUTHWPASS {
		z.Name = p.LastCommand
	}

	//Run the processor for the command issued by the client (if it exists) and return the result with the correct SMTP result prefix.
	if _, exists := Commands[strings.ToLower(z.Name)]; exists {
		p.LastCommand = z.Name
		res := Commands[strings.ToLower(z.Name)](p.Session, z)
		//Needed by the multiline EHLO response
		if res.Code == 0 {
			return res.Message
		}
		return res.String()
	}
	return "504 command not implemented"
}
