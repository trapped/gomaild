//Implements the STLS command.
package stls

import (
	"errors"
	"github.com/trapped/gomaild/cipher"
	"github.com/trapped/gomaild/config"
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/pop3/session"
	"log"
)

//Processes the STLS command.
func Process(session *Session, c Statement) (string, error) {
	if !config.Configuration.POP3.EnableSTLS {
		return "", errors.New("command not available")
	}

	//Check if the TLS cipher is not available (not initialized or not working)
	if !cipher.TLSAvailable {
		//Attempt to initialize it
		cipher.TLSLoadCertificate()
		//Check again to see if it is still not available
		if !cipher.TLSAvailable {
			return "", errors.New("error initializing the TLS cipher")
		}
	}

	if session.InTLS {
		return "", errors.New("already in TLS")
	}

	log.Println("POP3:", "STLS command issued by", session.RemoteEP)

	session.InTLS = true
	session.State = AUTHORIZATION //After issuing the STLS command, clients EHLO/HELO again, since the server can disclose more information
	return "ready to start TLS", nil
}
