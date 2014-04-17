package auth

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"github.com/trapped/gomaild/config"
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/smtp/reply"
	. "github.com/trapped/gomaild/processors/smtp/session"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

//Processes CRAM-MD5 authentication.
func CRAM_MD5(session *Session, c Statement) Reply {
	log.Println("SMTP:", "AUTH CRAM-MD5 (fragment) command issued by", session.RemoteEP)
	session.AuthMode = "cram-md5"
	shared := "<" + strconv.Itoa(os.Getpid()) + "." + strconv.Itoa(time.Now().Nanosecond()) + "@" + config.Configuration.ServerName + ">"
	if !config.Configuration.SMTP.EnableAUTH_CRAM_MD5 {
		return Reply{Code: 502, Message: "command not available"}
	}
	if len(c.Arguments) == 1 && session.AuthState != NONE && session.AuthState == AUTHWUSER {
		buf, err := base64.StdEncoding.DecodeString(c.Arguments[0])
		if err != nil {
			log.Println("SMTP:", "AUTH CRAM-MD5: Failed to decode from base64:", c.Arguments[0])
			return Reply{Code: 451, Message: "failed to decode from base64"}
		}

		//Parse the client's response to the challenge
		fields := strings.Split(string(buf), " ")
		if len(fields) < 2 {
			return Reply{Code: 501, Message: "wrong number of fields in the token"}
		}

		//Calculate the correct hash
		part1 := hmac.New(md5.New, []byte(config.Configuration.Accounts[fields[0]]))
		part1.Write([]byte(session.Shared))
		part2 := hex.EncodeToString(part1.Sum(nil))

		if part2 != fields[1] {
			session.Username = ""
			session.Password = ""
			session.AuthState = AUTHNONE
			session.AuthMode = ""
			return Reply{Code: 535, Message: "authentication failed"}
		}

		session.Username = fields[0]
		session.Password = config.Configuration.Accounts[fields[0]]
		session.AuthState = AUTHENTICATED
		log.Println("SMTP:", "AUTH LOGIN: Authentication successful for", session.RemoteEP)
		return Reply{Code: 235, Message: "authentication successful"}
	}

	if len(c.Arguments) != 2 {
		return Reply{Code: 501, Message: "wrong number of arguments"}
	}

	//Save the session state (waiting-for-input state and the shared that is sent to the client in the challenge)
	session.Shared = shared
	session.AuthState = AUTHWUSER
	//Send the challenge
	return Reply{Code: 334, Message: base64.StdEncoding.EncodeToString([]byte(shared))}
}
