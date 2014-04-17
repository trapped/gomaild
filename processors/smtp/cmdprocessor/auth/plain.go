package auth

import (
	"encoding/base64"
	"github.com/trapped/gomaild/config"
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/smtp/reply"
	. "github.com/trapped/gomaild/processors/smtp/session"
	"log"
	"strings"
)

//Processes PLAIN authentication.
func Plain(session *Session, c Statement) Reply {
	log.Println("SMTP:", "AUTH PLAIN (fragment) command issued by", session.RemoteEP)
	session.AuthMode = "plain"
	if !config.Configuration.SMTP.EnableAUTH_PLAIN {
		return Reply{Code: 502, Message: "command not available"}
	}
	//If waiting for data
	if len(c.Arguments) == 1 && session.AuthState != NONE && session.AuthState == AUTHWUSER {
		//Decode the data received
		buf, err := base64.StdEncoding.DecodeString(c.Arguments[2])
		if err != nil {
			log.Println("SMTP:", "AUTH PLAIN: Failed to decode from base64:", c.Arguments[0])
			return Reply{Code: 451, Message: "failed to decode from base64"}
		}

		//Parse the data
		fields := strings.Split(string(buf), "\x00")
		switch len(fields) {
		//Some clients omit the authentication-id field since it's useless
		case 2:
			session.Username = fields[0]
			session.Password = fields[1]
			break
		case 3:
			session.Username = fields[1]
			session.Password = fields[2]
			break
		default:
			return Reply{Code: 501, Message: "wrong number of fields in the token"}
		}

		//If the data doesn't match, reset the session
		if config.Configuration.Accounts[session.Username] != session.Password {
			session.Username = ""
			session.Password = ""
			session.AuthState = AUTHNONE
			session.AuthMode = ""
			return Reply{Code: 535, Message: "authentication failed"}
		}

		//Set the authentication state to "authenticated"
		session.AuthState = AUTHENTICATED
		log.Println("SMTP:", "AUTH LOGIN: Authentication successful for", session.RemoteEP)
		return Reply{Code: 235, Message: "authentication successful"}
	}

	switch len(c.Arguments) {
	//If no arguments beyond "AUTH PLAIN" has been provided
	case 2:
		//Set authentication state to "waiting for data"
		session.AuthState = AUTHWUSER
		return Reply{Code: 334, Message: ""}
	//If the client provided the PLAIN data already
	case 3:
		//Decode the data
		buf, err := base64.StdEncoding.DecodeString(c.Arguments[2])
		if err != nil {
			log.Println("SMTP:", "AUTH PLAIN: Failed to decode from base64:", c.Arguments[2])
			return Reply{Code: 451, Message: "failed to decode from base64"}
		}

		//Parse the data
		fields := strings.Split(string(buf), "\x00")
		switch len(fields) {
		//Some clients omit the authentication-id field since it's useless
		case 2:
			session.Username = fields[0]
			session.Password = fields[1]
			break
		case 3:
			session.Username = fields[1]
			session.Password = fields[2]
			break
		default:
			return Reply{Code: 501, Message: "wrong number of fields in the token"}
		}

		//If the data doesn't match, reset the session
		if config.Configuration.Accounts[session.Username] != session.Password {
			session.Username = ""
			session.Password = ""
			session.AuthState = AUTHNONE
			session.AuthMode = ""
			return Reply{Code: 535, Message: "authentication failed"}
		}

		//Set the authentication state to "authenticated"
		session.AuthState = AUTHENTICATED
		log.Println("SMTP:", "AUTH LOGIN: Authentication successful for", session.RemoteEP)
		return Reply{Code: 235, Message: "authentication successful"}
	}
	return Reply{Code: 451, Message: "processing error"}
}
