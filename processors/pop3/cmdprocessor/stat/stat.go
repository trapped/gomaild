package stat

import (
	"errors"
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/pop3/session"
	"log"
)

func Process(session *Session, c Command) (string, error) {
	if !session.Authenticated || session.User == "" || session.Password == "" {
		return "", errors.New("client not authenticated")
	}
}
