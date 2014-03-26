package capa

import (
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/pop3/session"
	"log"
)

func Process(session *Session, c Statement) (string, error) {
	log.Println("POP3:", "CAPA command issued by", session.RemoteEP, "with", session.Username)

	capabilities := "Capability list follows\r\n"
	capabilities += "TOP\r\n"
	capabilities += "USER\r\n"
	capabilities += "APOP\r\n"
	capabilities += "UIDL\r\n"
	capabilities += "PIPELINING\r\n"
	capabilities += "EXPIRE NEVER\r\n"
	capabilities += "RESP-CODES\r\n"
	capabilities += "IMPLEMENTATION gomaild_<http://github.com/trapped/gomaild>\r\n"
	capabilities += "."

	return capabilities, nil
}
