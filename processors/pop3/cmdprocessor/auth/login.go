package auth

import (
	"encoding/base64"
	"errors"
	"github.com/trapped/gomaild/config"
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/pop3/session"
	"log"
)

//Processes LOGIN authentication.
func Login(session *Session, c Statement) (string, error) {
	log.Println("POP3:", "AUTH LOGIN (fragment) command issued by", session.RemoteEP)
	session.AuthMode = "login"
	if !config.Configuration.POP3.EnableAUTH_LOGIN {
		return "", errors.New("command not available")
	}
	//If waiting for data, only one argument provided and not authenticated already
	if len(c.Arguments) == 1 && session.AuthState != AUTHNONE && session.AuthState != AUTHENTICATED {
		if session.AuthState == AUTHWUSER {
			//Decode the username the client has provided
			buf, err := base64.StdEncoding.DecodeString(c.Arguments[0])
			if err != nil {
				log.Println("POP3:", "AUTH LOGIN: Failed to decode from base64:", c.Arguments[0])
				return "", errors.New("failed to decode from base64")
			}

			//Set the state to "Waiting for password"
			session.AuthState = AUTHWPASS
			//Save the username in the session
			session.Username = string(buf)
			log.Println("POP3:", "AUTH LOGIN: Username", session.Username, "received from", session.RemoteEP)
			return base64.StdEncoding.EncodeToString([]byte("Password:")), nil
		} else if session.AuthState == AUTHWPASS {
			//Decode the password the client has provided
			buf, err := base64.StdEncoding.DecodeString(c.Arguments[0])
			if err != nil {
				log.Println("POP3:", "AUTH LOGIN: Failed to decode from base64:", c.Arguments[0])
				return "", errors.New("failed to decode from base64")
			}

			//Save the password in the session
			session.Password = string(buf)

			//If the password is incorrect, reset authentication state
			if config.Configuration.Accounts[session.Username] != session.Password {
				session.Username = ""
				session.Password = ""
				session.AuthState = AUTHNONE
				session.AuthMode = ""
				log.Println("POP3:", "AUTH LOGIN: Authentication failed for", session.RemoteEP)
				return "", errors.New("authentication failed")
			}

			//Set authentication state to "authenticated"
			session.AuthState = AUTHENTICATED
			session.State = TRANSACTION
			session.Authenticated = true
			log.Println("POP3:", "AUTH LOGIN: Authentication successful for", session.RemoteEP)
			return "authentication successful", nil
		}
	}

	switch len(c.Arguments) {
	//If only "AUTH LOGIN" has been sent
	case 2:
		//Set state to "waiting for user"
		session.AuthState = AUTHWUSER
		return base64.StdEncoding.EncodeToString([]byte("Username:")), nil
	//If "AUTH LOGIN username" has been sent
	case 3:
		//Decode the username
		buf, err := base64.StdEncoding.DecodeString(c.Arguments[2])
		if err != nil {
			log.Println("POP3:", "AUTH LOGIN: Failed to decode from base64:", c.Arguments[2])
			return "", errors.New("failed to decode from base64")
		}
		//Set state to "waiting for password"
		session.AuthState = AUTHWPASS
		//Save the username in the session
		session.Username = string(buf)
		log.Println("POP3:", "AUTH LOGIN: Username", session.Username, "received from", session.RemoteEP)
		return base64.StdEncoding.EncodeToString([]byte("Password:")), nil
	}
	//This should never happen
	return "", errors.New("processing error")
}
