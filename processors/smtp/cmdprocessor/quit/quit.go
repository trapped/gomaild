package quit

import (
	"github.com/trapped/gomaild/config"
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/smtp/message"
	. "github.com/trapped/gomaild/processors/smtp/reply"
	. "github.com/trapped/gomaild/processors/smtp/session"
	"log"
)

func Process(session *Session, c Statement) Reply {
	if len(c.Arguments) != 1 {
		return Reply{Code: 501, Message: "too many arguments"}
	}

	log.Println("SMTP: QUIT command issued by", session.RemoteEP)

	if len(session.Received) != 0 {
		for _, v := range session.Received {
			err := Store(session, *v.(*Message))
			if err != nil {
				log.Println("SMTP:", "Error storing message:", err)
			}
		}
	}

	session.Quitted = true

	return Reply{Code: 221, Message: config.Configuration.SMTP.EndGreeting}
}
