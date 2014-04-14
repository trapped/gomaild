package starttls

import (
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/smtp/reply"
	. "github.com/trapped/gomaild/processors/smtp/session"
	"log"
)

func Process(session *Session, c Statement) Reply {
	if !config.Configuration.SMTP.EnableSTARTTLS {
		return Reply{Code: 502, Message: "command not available"}
	}
	if session.InTLS {
		return Reply{Code: 454, Message: "already in TLS"}
	}

	log.Println("SMTP:", "STARTTLS command issued by", session.RemoteEP)

	session.InTLS = true
	return Reply{Code: 220, Message: "ready to start TLS"}
}
