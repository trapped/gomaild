package rset

import (
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/smtp/reply"
	. "github.com/trapped/gomaild/processors/smtp/session"
	"log"
)

func Process(session *Session, c Statement) Reply {
	if len(c.Arguments) != 1 {
		return Reply{Code: 501, Message: "too many arguments"}
	}

	log.Println("SMTP:", "RSET command issued by", session.RemoteEP, "with", session.Identity)

	session.Received = []interface{}{}
	session.Shared = ""
	session.State = NONE

	return Reply{Code: 250, Message: "session has been reset"}
}
