//Package auth implements SMTP authentication, such as PLAIN, LOGIN, and CRAM-MD5 SASL methods.
package auth

import (
	"github.com/trapped/gomaild/config"
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/smtp/reply"
	. "github.com/trapped/gomaild/processors/smtp/session"
	"strings"
)

//Processes authentication commands and data.
func Process(session *Session, c Statement) Reply {
	if !config.Configuration.SMTP.EnableAUTH {
		return Reply{Code: 502, Message: "command not available"}
	}
	if session.State != IDENTIFIED && session.State != COMPOSITION {
		return Reply{Code: 503, Message: "identify first"}
	}
	if session.Authenticated {
		return Reply{Code: 503, Message: "already authenticated"}
	}
	if len(c.Arguments) < 2 && session.AuthState == AUTHNONE {
		return Reply{Code: 501, Message: "wrong number of arguments"}
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
		return Reply{Code: 502, Message: "authentication method not implemented"}
	}
	return Reply{Code: 502, Message: "authentication method not implemented"}
}
