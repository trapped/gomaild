package helo

import (
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/smtp/error"
	. "github.com/trapped/gomaild/processors/smtp/reply"
	. "github.com/trapped/gomaild/processors/smtp/session"
	"log"
)

func Process(session *Session, c Statement) (Reply, Error) {
	if session.State != NONE {
		return Reply{}, Error{Code: 503, Message: "wrong session state"}
	}
	if len(c.Arguments) != 2 {
		return Reply{}, Error{Code: 501, Message: "wrong number of arguments"}
	}

	session.Identity = c.Arguments[1]

	log.Println("SMTP:", "HELO command issued by", session.RemoteEP, "with", session.Identity)

	greeting := "250 greetings, " + session.Identity + "\r\n"

    session.State = IDENTIFICATED

	return Reply{Message: greeting}, Error{}
}
