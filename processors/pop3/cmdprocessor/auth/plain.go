package auth

import (
	"encoding/base64"
	"errors"
	"github.com/trapped/gomaild/config"
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/pop3/session"
	"log"
	"strings"
)

//Processes PLAIN authentication.
func Plain(session *Session, c Statement) (string, error) {
	log.Println("POP3:", "AUTH PLAIN (fragment) command issued by", session.RemoteEP)
	session.AuthMode = "plain"
	if !config.Configuration.POP3.EnableAUTH_PLAIN {
		return "", errors.New("command not available")
	}
	//If waiting for data
	if len(c.Arguments) == 1 && session.AuthState != AUTHNONE && session.AuthState == AUTHWUSER {
		//Decode the data received
		buf, err := base64.StdEncoding.DecodeString(c.Arguments[2])
		if err != nil {
			log.Println("POP3:", "AUTH PLAIN: Failed to decode from base64:", c.Arguments[0])
			return "", errors.New("failed to decode from base64")
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
			return "", errors.New("wrong number of fields in the token")
		}

		//If the data doesn't match, reset the session
		if config.Configuration.Accounts[session.Username] != session.Password {
			session.Username = ""
			session.Password = ""
			session.AuthState = AUTHNONE
			session.AuthMode = ""
			return "", errors.New("authentication failed")
		}

		//Set the authentication state to "authenticated"
		session.AuthState = AUTHENTICATED
		session.State = TRANSACTION
		session.Authenticated = true
		log.Println("POP3:", "AUTH LOGIN: Authentication successful for", session.RemoteEP)
		return "authentication successful", nil
	}

	switch len(c.Arguments) {
	//If no arguments beyond "AUTH PLAIN" has been provided
	case 2:
		//Set authentication state to "waiting for data"
		session.AuthState = AUTHWUSER
		return "", nil
	//If the client provided the PLAIN data already
	case 3:
		//Decode the data
		buf, err := base64.StdEncoding.DecodeString(c.Arguments[2])
		if err != nil {
			log.Println("POP3:", "AUTH PLAIN: Failed to decode from base64:", c.Arguments[2])
			return "", errors.New("failed to decode from base64")
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
			return "", errors.New("wrong number of fields in the token")
		}

		//If the data doesn't match, reset the session
		if config.Configuration.Accounts[session.Username] != session.Password {
			session.Username = ""
			session.Password = ""
			session.AuthState = AUTHNONE
			session.AuthMode = ""
			return "", errors.New("authentication failed")
		}

		//Set the authentication state to "authenticated"
		session.AuthState = AUTHENTICATED
		session.Authenticated = true
		session.State = TRANSACTION
		log.Println("POP3:", "AUTH LOGIN: Authentication successful for", session.RemoteEP)
		return "authentication successful", nil
	}
	return "", errors.New("processing error")
}
