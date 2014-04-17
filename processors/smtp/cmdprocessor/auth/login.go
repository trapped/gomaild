package auth

import (
	"encoding/base64"
	"github.com/trapped/gomaild/config"
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/smtp/reply"
	. "github.com/trapped/gomaild/processors/smtp/session"
	"log"
)

//Processes LOGIN authentication.
func Login(session *Session, c Statement) Reply {
	log.Println("SMTP:", "AUTH LOGIN (fragment) command issued by", session.RemoteEP)
	session.AuthMode = "login"
	if !config.Configuration.SMTP.EnableAUTH_LOGIN {
		return Reply{Code: 502, Message: "command not available"}
	}
	//If waiting for data, only one argument provided and not authenticated already
	if len(c.Arguments) == 1 && session.AuthState != NONE && session.AuthState != AUTHENTICATED {
		if session.AuthState == AUTHWUSER {
			//Decode the username the client has provided
			buf, err := base64.StdEncoding.DecodeString(c.Arguments[0])
			if err != nil {
				log.Println("SMTP:", "AUTH LOGIN: Failed to decode from base64:", c.Arguments[0])
				return Reply{Code: 451, Message: "failed to decode from base64"}
			}

			//Set the state to "Waiting for password"
			session.AuthState = AUTHWPASS
			//Save the username in the session
			session.Username = string(buf)
			log.Println("SMTP:", "AUTH LOGIN: Username", session.Username, "received from", session.RemoteEP)
			return Reply{Code: 334, Message: base64.StdEncoding.EncodeToString([]byte("Password:"))}
		} else if session.AuthState == AUTHWPASS {
			//Decode the password the client has provided
			buf, err := base64.StdEncoding.DecodeString(c.Arguments[0])
			if err != nil {
				log.Println("SMTP:", "AUTH LOGIN: Failed to decode from base64:", c.Arguments[0])
				return Reply{Code: 451, Message: "failed to decode from base64"}
			}

			//Save the password in the session
			session.Password = string(buf)

			//If the password is incorrect, reset authentication state
			if config.Configuration.Accounts[session.Username] != session.Password {
				session.Username = ""
				session.Password = ""
				session.AuthState = AUTHNONE
				session.AuthMode = ""
				log.Println("SMTP:", "AUTH LOGIN: Authentication failed for", session.RemoteEP)
				return Reply{Code: 535, Message: "authentication failed"}
			}

			//Set authentication state to "authenticated"
			session.AuthState = AUTHENTICATED
			log.Println("SMTP:", "AUTH LOGIN: Authentication successful for", session.RemoteEP)
			return Reply{Code: 235, Message: "authentication successful"}
		}
	}

	switch len(c.Arguments) {
	//If only "AUTH LOGIN" has been sent
	case 2:
		//Set state to "waiting for user"
		session.AuthState = AUTHWUSER
		return Reply{Code: 334, Message: base64.StdEncoding.EncodeToString([]byte("Username:"))}
	//If "AUTH LOGIN username" has been sent
	case 3:
		//Decode the username
		buf, err := base64.StdEncoding.DecodeString(c.Arguments[2])
		if err != nil {
			log.Println("SMTP:", "AUTH LOGIN: Failed to decode from base64:", c.Arguments[2])
			return Reply{Code: 451, Message: "failed to decode from base64"}
		}
		//Set state to "waiting for password"
		session.AuthState = AUTHWPASS
		//Save the username in the session
		session.Username = string(buf)
		log.Println("SMTP:", "AUTH LOGIN: Username", session.Username, "received from", session.RemoteEP)
		return Reply{Code: 334, Message: base64.StdEncoding.EncodeToString([]byte("Password:"))}
	}
	//This should never happen
	return Reply{Code: 451, Message: "processing error"}
}
