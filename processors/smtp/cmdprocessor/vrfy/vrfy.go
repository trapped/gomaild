package vrfy

import (
	"github.com/trapped/gomaild/mailboxes"
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/smtp/reply"
	. "github.com/trapped/gomaild/processors/smtp/session"
	"log"
	"regexp"
)

func Process(session *Session, c Statement) Reply {
	if session.State != IDENTIFICATED {
		return Reply{Code: 503, Message: "wrong session state"}
	}
	if len(c.Arguments) < 2 {
		return Reply{Code: 501, Message: "not enough arguments"}
	}

	log.Println("SMTP:", "VRFY command issued by", session.RemoteEP, "with", session.Identity)

	data := c.Arguments[1]
	regex, err := regexp.Compile("(i?)[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a0-9])?")
	if err != nil {
		return Reply{Code: 451, Message: "processing error while parsing the regex"}
	}
	if !regex.MatchString(data) {
		return Reply{Code: 553, Message: "invalid address"}
	}

	_, err = mailboxes.GetUser(c.Arguments[1])
	if err != nil {
		return Reply{Code: 251, Message: "user might exist"}
	}

	return Reply{Code: 251, Message: "user might exist"}
}
