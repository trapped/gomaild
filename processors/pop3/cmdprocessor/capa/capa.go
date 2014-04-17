//Implements the CAPA command.
package capa

import (
	"github.com/trapped/gomaild/config"
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/pop3/session"
	"log"
	"strings"
)

//Processes the CAPA command.
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
	if config.Configuration.POP3.EnableSTLS {
		capabilities += "STLS\r\n"
	}
	if config.Configuration.POP3.EnableAUTH {
		auths := []string{}
		if config.Configuration.POP3.EnableAUTH_LOGIN {
			auths = append(auths, "LOGIN")
		}
		if config.Configuration.POP3.EnableAUTH_PLAIN {
			auths = append(auths, "PLAIN")
		}
		if config.Configuration.POP3.EnableAUTH_CRAM_MD5 {
			auths = append(auths, "CRAM-MD5")
		}
		if len(auths) > 0 {
			capabilities += "AUTH " + strings.Join(auths, " ") + "\r\n"
		}
	}
	capabilities += "IMPLEMENTATION gomaild_<http://github.com/trapped/gomaild>\r\n"
	capabilities += "."

	return capabilities, nil
}
