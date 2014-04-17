//Implements the NOOP command.
package noop

import (
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/pop3/session"
	"log"
)

//Processes the NOOP command.
func Process(session *Session, c Statement) (string, error) {
	log.Println("POP3:", "NOOP command issued by", session.RemoteEP, "with", session.Username)
	return "", nil
}
