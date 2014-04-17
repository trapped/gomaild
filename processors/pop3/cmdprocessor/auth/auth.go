//Package auth implements POP3 authentication, such as PLAIN, LOGIN, and CRAM-MD5 SASL methods.
package auth

import (
	"errors"
	"github.com/trapped/gomaild/config"
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/pop3/session"
	"strings"
)

//Processes authentication commands and data.
func Process(session *Session, c Statement) (string, error) {
	if !config.Configuration.POP3.EnableAUTH {
		return "", errors.New("command not available")
	}
	if session.State != AUTHORIZATION {
		return "", errors.New("wrong session state")
	}
	if session.Authenticated {
		return "", errors.New("already authenticated")
	}
	if len(c.Arguments) < 2 && session.AuthState == AUTHNONE {
		return "", errors.New("wrong number of arguments")
	}

	mode := ""

	//Check if the authentication process has already started and the received data has to be processed by the last used AuthMode
	if len(c.Arguments) < 2 && session.AuthMode != "" {
		mode = session.AuthMode
	} else {
		mode = c.Arguments[1]
	}

	switch strings.ToLower(mode) {
	case "plain":
		return Plain(session, c)
		break
	case "login":
		return Login(session, c)
		break
	case "cram-md5":
		return CRAM_MD5(session, c)
		break
	/*case "digest-md5":
	return Digest_MD5(session, c)
	break*/
	default:
		return "", errors.New("authentication method not implemented")
	}
	return "", errors.New("authentication method not implemented")
}
