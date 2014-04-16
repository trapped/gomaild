package ehlo

import (
	"github.com/trapped/gomaild/config"
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/smtp/reply"
	. "github.com/trapped/gomaild/processors/smtp/session"
	"log"
	"strings"
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
	capabilities += "250-8BITMIME\r\n" //I'm not really sure that Go natively supports UTF8/Unicode text, but in the end it's all binary data, and I'm just dumping it to a file, right?
	if !session.InTLS && config.Configuration.SMTP.EnableSTARTTLS {
		capabilities += "250-STARTTLS\r\n"
	}
	if config.Configuration.SMTP.EnableAUTH {
		auths := []string{}
		if config.Configuration.SMTP.EnableAUTH_LOGIN {
			auths = append(auths, "LOGIN")
		}
		if config.Configuration.SMTP.EnableAUTH_PLAIN {
			auths = append(auths, "PLAIN")
		}
		if len(auths) != 0 {
			capabilities += "250-AUTH " + strings.Join(auths, " ") + "\r\n"
		}
	}
	capabilities += "250 PIPELINING"

	session.State = IDENTIFIED

	return Reply{Message: capabilities}
}
