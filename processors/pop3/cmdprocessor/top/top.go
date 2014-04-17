//Implements the TOP command.
package top

import (
	"errors"
	. "github.com/trapped/gomaild/parsers/textual"
	"github.com/trapped/gomaild/processors/pop3/message"
	. "github.com/trapped/gomaild/processors/pop3/session"
	"log"
	"strconv"
	"strings"
)

//Processes the TOP command.
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
	if len(c.Arguments) > 3 {
		errorslice = append(errorslice, "too many arguments")
	}
	if len(c.Arguments) < 3 {
		errorslice = append(errorslice, "not enough arguments")
	}

	if len(errorslice) != 0 {
		goto returnerror
	}

	log.Println("POP3:", "TOP command issued by", session.RemoteEP, "with", session.Username)

	messages := message.Index(session)

	for _, v := range messages {
		if strconv.Itoa(v.ID) == c.Arguments[1] {
			result = "top of message follows\r\n"

			headers, err := message.Headers(v)
			if err != nil {
				errorslice = append(errorslice, "error reading the message headers")
				goto returnerror
			}

			result += headers + "\r\n"

			body, err := message.Body(v)
			if err != nil {
				errorslice = append(errorslice, "error reading the message body")
			}

			times, err := strconv.Atoi(c.Arguments[2])
			if err != nil {
				errorslice = append(errorslice, "error parsing the second TOP argument")
			}

			if times > 0 {
				result += "\r\n"

				bodylines := strings.Split(body, "\r\n")

				if times > len(bodylines) {
					result += strings.Join(bodylines, "\r\n")
				} else {
					log.Println(times)
					result += strings.Join(bodylines[:times], "\r\n")
				}
			}

			result += "\r\n."

			session.Retrieved = append(session.Retrieved, v)
			break
		}
	}

	if result == "" {
		errorslice = append(errorslice, "no such message; "+strconv.Itoa(len(messages))+" messages in maildrop")
		goto returnerror
	}

	return result, nil
}
