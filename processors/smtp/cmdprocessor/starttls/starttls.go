package starttls

import (
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/smtp/reply"
	. "github.com/trapped/gomaild/processors/smtp/session"
	"log"
)

func Process(session *Session, c Statement) Reply {
	if session.InTLS {
		return Reply{Code: 454, Message: "already in TLS"}
	}
	session.State = NONE
	session.Identity = ""
	return Reply{Code: 220, Message: "ready to start TLS"}
}
