package user

import (
	"errors"
	"github.com/trapped/gomaild/config"
	"github.com/trapped/gomaild/mailboxes"
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/pop3/session"
	"log"
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
	if !config.Configuration.POP3.EnableUSER {
		errorslice = append(errorslice, "command not available")
		goto returnerror
	}
	if session.State != AUTHORIZATION {
		errorslice = append(errorslice, "wrong session state")
	}
	if session.Authenticated {
		errorslice = append(errorslice, "already authenticated")
	}
	if session.Username != "" {
		errorslice = append(errorslice, "user already set, use command PASS")
	}
	if len(c.Arguments) == 1 {
		errorslice = append(errorslice, "username can't be empty")
	}
	if len(c.Arguments) > 2 {
		errorslice = append(errorslice, "too many arguments")
	}

	if len(errorslice) != 0 {
		goto returnerror
	}

	log.Println("POP3:", "USER command issued by", session.RemoteEP, "with", c.Arguments[1])

	if !config.Configuration.POP3.SecureUSER {
		_, erra := mailboxes.GetUser(c.Arguments[1])
		if erra != nil {
			errorslice = append(errorslice, "no such user")
			goto returnerror
		}
	}

	session.Username = c.Arguments[1]
	result += "user might exist"

	return result, nil
}
