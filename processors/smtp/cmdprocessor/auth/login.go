package auth

import (
	"encoding/base64"
	"github.com/trapped/gomaild/config"
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/smtp/reply"
	. "github.com/trapped/gomaild/processors/smtp/session"
	"log"
)

func Login(session *Session, c Statement) Reply {
	log.Println("SMTP:", "AUTH LOGIN (fragment) command issued by", session.RemoteEP)
	session.AuthMode = "login"
	if !config.Configuration.SMTP.EnableAUTH_LOGIN {
		return Reply{Code: 502, Message: "command not available"}
	}
	if len(c.Arguments) == 1 && session.AuthState != NONE && session.AuthState != AUTHENTICATED {
		if session.AuthState == AUTHWUSER {
			buf, err := base64.StdEncoding.DecodeString(c.Arguments[0])
			if err != nil {
				log.Println("SMTP:", "AUTH LOGIN: Failed to decode from base64:", c.Arguments[0])
				return Reply{Code: 451, Message: "failed to decode from base64"}
			}
			session.AuthState = AUTHWPASS
			session.Username = string(buf)
			log.Println("SMTP:", "AUTH LOGIN: Username", session.Username, "received from", session.RemoteEP)
			return Reply{Code: 334, Message: base64.StdEncoding.EncodeToString([]byte("Password:"))}
		} else if session.AuthState == AUTHWPASS {
			buf, err := base64.StdEncoding.DecodeString(c.Arguments[0])
			if err != nil {
				log.Println("SMTP:", "AUTH LOGIN: Failed to decode from base64:", c.Arguments[0])
				return Reply{Code: 451, Message: "failed to decode from base64"}
			}
			session.Password = string(buf)

			if config.Configuration.Accounts[session.Username] != session.Password {
				session.Username = ""
				session.Password = ""
				session.AuthState = AUTHNONE
				session.AuthMode = ""
				log.Println("SMTP:", "AUTH LOGIN: Authentication failed for", session.RemoteEP)
				return Reply{Code: 535, Message: "authentication failed"}
			}

			session.AuthState = AUTHENTICATED
			session.AuthMode = ""
			log.Println("SMTP:", "AUTH LOGIN: Authentication successful for", session.RemoteEP)
			return Reply{Code: 235, Message: "authentication successful"}
		}
	}

	switch len(c.Arguments) {
	case 2:
		session.AuthState = AUTHWUSER
		return Reply{Code: 334, Message: base64.StdEncoding.EncodeToString([]byte("Username:"))}
	case 3:
		buf, err := base64.StdEncoding.DecodeString(c.Arguments[2])
		if err != nil {
			log.Println("SMTP:", "AUTH LOGIN: Failed to decode from base64:", c.Arguments[2])
			return Reply{Code: 451, Message: "failed to decode from base64"}
		}
		session.AuthState = AUTHWPASS
		session.Username = string(buf)
		log.Println("SMTP:", "AUTH LOGIN: Username", session.Username, "received from", session.RemoteEP)
		return Reply{Code: 334, Message: base64.StdEncoding.EncodeToString([]byte("Password:"))}
	}
	return Reply{Code: 451, Message: "processing error"}
}
