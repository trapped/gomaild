package noop

import (
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/smtp/reply"
	. "github.com/trapped/gomaild/processors/smtp/session"
	"log"
)

func Process(session *Session, c Statement) Reply {
	log.Println("SMTP:", "NOOP command issued by", session.RemoteEP, "with", session.Identity)
	return Reply{Code: 250, Message: "no operation performed"}
}
