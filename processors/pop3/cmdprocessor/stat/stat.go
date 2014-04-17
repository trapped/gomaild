//Implements the STAT command.
package stat

import (
	"errors"
	"github.com/trapped/gomaild/mailboxes"
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/pop3/session"
	"log"
	"strconv"
	"strings"
)

//Processes the STAT command.
func Process(session *Session, c Statement) (string, error) {
	errorslice := []string{}
	result := ""
	goto checks

returnerror:
	if len(errorslice) != 0 {
		result = strings.Join(errorslice, ", ")
		return "", errors.New(result)
	}

checks:
	if !session.Authenticated {
		errorslice = append(errorslice, "not authenticated")
	}
	if session.Username == "" {
		errorslice = append(errorslice, "session username can't be empty")
	}
	if session.Password == "" {
		errorslice = append(errorslice, "session password can't be empty")
	}
	if session.State != TRANSACTION {
		errorslice = append(errorslice, "wrong state")
	}
	if len(c.Arguments) != 1 {
		errorslice = append(errorslice, "too many arguments")
	}

	if len(errorslice) != 0 {
		goto returnerror
	}

	log.Println("POP3:", "STAT command issued by", session.RemoteEP, "with", session.Username)

	count, octets := mailboxes.Stat(session.Username, false)

	return strconv.Itoa(count) + " " + strconv.Itoa(octets), nil
}
