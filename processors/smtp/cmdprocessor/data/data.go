package data

import (
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/smtp/message"
	. "github.com/trapped/gomaild/processors/smtp/reply"
	. "github.com/trapped/gomaild/processors/smtp/session"
	"log"
)

func Process(session *Session, c Statement) Reply {
	if session.State != RECAPITATION && session.State != COMPOSITION {
		return Reply{Code: 503, Message: "wrong session state"}
	}
	if len(c.Arguments) != 1 && session.State != COMPOSITION {
		return Reply{Code: 501, Message: "too many arguments"}
	}

	log.Println("SMTP:", "DATA command issued by", session.RemoteEP, "with", session.Identity)

	if session.State == RECAPITATION {
		session.State = COMPOSITION
		return Reply{Code: 354, Message: "start mail input, end with <CRLF.CRLF>"}
	}

	if c.Raw == ".\r\n" {
		session.State = IDENTIFICATED
		return Reply{Code: 250, Message: "message accepted and queued"}
	}

	session.Received[len(session.Received)-1].(*Message).Text += c.Raw

	return Reply{}
}
