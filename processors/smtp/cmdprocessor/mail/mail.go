package mail

import (
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/smtp/message"
	. "github.com/trapped/gomaild/processors/smtp/reply"
	. "github.com/trapped/gomaild/processors/smtp/session"
	"log"
	"regexp"
	"strings"
)

func Process(session *Session, c Statement) Reply {
	if session.State != IDENTIFICATED {
		return Reply{Code: 503, Message: "wrong session state"}
	}
	if len(c.Arguments) < 2 {
		return Reply{Code: 501, Message: "not enough arguments"}
	}

	log.Println("SMTP:", "MAIL command issued by", session.RemoteEP, "with", session.Identity)

	tempdata := make(map[string]string, 0)

	for i, v := range c.Arguments {
		switch strings.ToLower(strings.TrimSuffix(v, ":")) {
		case "from":
			sender := c.Arguments[i+1]
			regex, err := regexp.Compile("(i?)[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a0-9])?")
			if err != nil {
				return Reply{Code: 451, Message: "processing error while parsing the regex"}
			}
			if !regex.MatchString(sender) {
				return Reply{Code: 553, Message: "invalid sender address"}
			}
			if tempdata["sender"] != "" {
				return Reply{Code: 501, Message: "sender address cannot be set twice"}
			}
			tempdata["sender"] = sender
			break
		}
	}

	if tempdata["sender"] == "" {
		return Reply{Code: 501, Message: "sender not specified"}
	}

	session.State = RECAPITATION
	session.Received = append(session.Received, &Message{Sender: tempdata["sender"]})

	return Reply{Code: 250, Message: "sender ok"}
}
