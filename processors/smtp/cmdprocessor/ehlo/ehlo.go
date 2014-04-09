package ehlo

import (
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/smtp/reply"
	. "github.com/trapped/gomaild/processors/smtp/session"
	"log"
)

func Process(session *Session, c Statement) Reply {
	if session.State != NONE {
		return Reply{Code: 503, Message: "wrong session state"}
	}
	if len(c.Arguments) != 2 {
		return Reply{Code: 501, Message: "wrong number of arguments"}
	}

	session.Identity = c.Arguments[1]

	log.Println("SMTP:", "EHLO command issued by", session.RemoteEP, "with", session.Identity)

	capabilities := "250-greetings, " + session.Identity + "\r\n"
	capabilities += "250-8BITMIME\r\n"
	/*if !session.InTLS {
	    capabilities += "250-STARTTLS\r\n"
	}*/
	capabilities += "250 PIPELINING"

	session.State = IDENTIFICATED

	return Reply{Message: capabilities}
}
