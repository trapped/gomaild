package quit

import (
	"errors"
	"github.com/trapped/gomaild/config"
	"github.com/trapped/gomaild/mailboxes"
	. "github.com/trapped/gomaild/parsers/textual"
	"github.com/trapped/gomaild/processors/pop3/message"
	. "github.com/trapped/gomaild/processors/pop3/session"
	"log"
	"strconv"
	"strings"
)

func Process(session *Session, c Statement) (string, error) {
	errorslice := []string{}
	result := ""
	goto checks

returnerror:
	if len(errorslice) != 0 {
		session.Username = ""
		session.Password = ""
		result = strings.Join(errorslice, ", ")
		return "", errors.New(result)
	}

checks:
	if len(c.Arguments) != 1 {
		errorslice = append(errorslice, "too many arguments")
	}

	if len(errorslice) != 0 {
		goto returnerror
	}

	log.Println("POP3: QUIT command issued by", session.RemoteEP)

	if session.State == TRANSACTION {
		session.State = UPDATE
	}

	if session.State == UPDATE {
		currentmessages := message.Index(session)

		for i, _ := range session.Retrieved {
			for f, _ := range currentmessages {
				if session.Retrieved[i].(message.Message).ID == currentmessages[f].ID {
					message.MoveMessage(session, currentmessages[f], "read")
				}
			}
		}

		for i, _ := range session.Deleted {
			for f, _ := range currentmessages {
				if session.Deleted[i].(message.Message).ID == currentmessages[f].ID {
					if config.Configuration.POP3.FakeDELE {
						err := message.MoveMessage(session, currentmessages[f], "deleted")
						if err != nil {
							log.Println("POP3:", "Error fake-deleting message", err)
						}
					} else {
						err := message.DeleteMessage(currentmessages[f])
						if err != nil {
							log.Println("POP3:", "Error deleting message:", err)
						}
					}
				}
			}
		}
	}

	session.Quitted = true

	count, _ := mailboxes.Stat(session.Username, false)

	ex := config.Configuration.POP3.EndGreeting
	if session.Authenticated {
		if ex != "" {
			ex += " "
		}
		ex += "(" + strconv.Itoa(count) + " messages left)"
	}

	return ex, nil
}
