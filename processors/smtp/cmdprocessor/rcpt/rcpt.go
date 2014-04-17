//Implements the RCPT command.
package rcpt

import (
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/smtp/message"
	. "github.com/trapped/gomaild/processors/smtp/reply"
	. "github.com/trapped/gomaild/processors/smtp/session"
	"log"
	"regexp"
	"strings"
)

//Processes the RCPT command.
func Process(session *Session, c Statement) Reply {
	if session.State != RECAPITATION {
		return Reply{Code: 503, Message: "wrong session state"}
	}
	if len(c.Arguments) < 2 {
		return Reply{Code: 501, Message: "not enough arguments"}
	}

	log.Println("SMTP:", "RCPT command issued by", session.RemoteEP, "with", session.Identity)

	tempdata := make(map[string]string, 0)

	for i, v := range c.Arguments {
		switch strings.ToLower(strings.TrimSuffix(v, ":")) {
		case "to":
			recipient := c.Arguments[i+1]
			regex, err := regexp.Compile("(i?)[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a0-9])?")
			if err != nil {
				return Reply{Code: 451, Message: "processing error while parsing the regex"}
			}
			if !regex.MatchString(recipient) {
				return Reply{Code: 553, Message: "invalid recipient address"}
			}
			if tempdata["recipient"] != "" {
				return Reply{Code: 501, Message: "recipient address cannot be set twice"}
			}
			tempdata["recipient"] = recipient
			break
		}
	}

	if tempdata["recipient"] == "" {
		return Reply{Code: 501, Message: "recipient not specified"}
	}

	session.Received[len(session.Received)-1].(*Message).Recipients = append(session.Received[len(session.Received)-1].(*Message).Recipients, tempdata["recipient"])

	return Reply{Code: 250, Message: "recipient ok"}
}
