package user

import (
	"errors"
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/pop3/session"
	"log"
)

func Process(session *Session, c Command) (string, error) {
	if !session.Authenticated && session.User == "" && session.Password == "" {
		success := false
		if success {
			session.User = c.Arguments[1]
			return "existing user", nil
		} else {
			return "", errors.New("no such user")
		}
	}
}
