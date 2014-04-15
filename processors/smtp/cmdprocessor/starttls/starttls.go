package starttls

import (
	"github.com/trapped/gomaild/cipher"
	"github.com/trapped/gomaild/config"
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/smtp/reply"
	. "github.com/trapped/gomaild/processors/smtp/session"
	"log"
)

func Process(session *Session, c Statement) Reply {
	if !config.Configuration.SMTP.EnableSTARTTLS {
		return Reply{Code: 502, Message: "command not available"}
	}

	//Check if the TLS cipher is not available (not initialized or not working)
	if !cipher.TLSAvailable {
		//Attempt to initialize it
		cipher.TLSLoadCertificate()
		//Check again to see if it is still not available
		if !cipher.TLSAvailable {
			return Reply{Code: 451, Message: "error initializing the TLS cipher"}
		}
	}

	if session.InTLS {
		return Reply{Code: 454, Message: "already in TLS"}
	}

	log.Println("SMTP:", "STARTTLS command issued by", session.RemoteEP)

	session.InTLS = true
	return Reply{Code: 220, Message: "ready to start TLS"}
}
