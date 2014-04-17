//Implements the LIST command.
package list

import (
	"errors"
	"github.com/trapped/gomaild/mailboxes"
	. "github.com/trapped/gomaild/parsers/textual"
	"github.com/trapped/gomaild/processors/pop3/message"
	. "github.com/trapped/gomaild/processors/pop3/session"
	"log"
	"strconv"
	"strings"
)

//Processes the LIST command.
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
	if session.State != TRANSACTION {
		errorslice = append(errorslice, "wrong session state")
	}
	if !session.Authenticated {
		errorslice = append(errorslice, "not authenticated")
	}
	if session.Username == "" {
		errorslice = append(errorslice, "user can't be empty")
	}
	if len(c.Arguments) > 2 {
		errorslice = append(errorslice, "too many arguments")
	}

	if len(errorslice) != 0 {
		goto returnerror
	}

	log.Println("POP3:", "LIST command issued by", session.RemoteEP, "with", session.Username)

	messages := message.Index(session)
	if len(c.Arguments) == 1 {
		count, octets := mailboxes.Stat(session.Username, false)
		result = strconv.Itoa(count) + " messages (" + strconv.Itoa(octets) + " octets)\r\n"
		for i := 0; i < len(messages); i++ {
			result += strconv.Itoa(messages[i].ID) + " " + strconv.Itoa(int(messages[i].File.Size())) + "\r\n"
		}
		result += "."
	} else {
		for _, v := range messages {
			if strconv.Itoa(v.ID) == c.Arguments[1] {
				result = strconv.Itoa(v.ID) + " " + strconv.Itoa(int(v.File.Size()))
				break
			}
		}
		if result == "" {
			errorslice = append(errorslice, "no such message; "+strconv.Itoa(len(messages))+" messages in maildrop")
			goto returnerror
		}
	}

	return result, nil
}
