package apop

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/trapped/gomaild/mailboxes"
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/pop3/session"
	"log"
	"strconv"
	"strings"
)

func Process(session *Session, c Statement) (string, error) {
	errorslice := []string{}
	result := ""
	goto checks

returnerror:
	if len(errorslice) != 0 {
		session.Username = ""
		session.Password = ""
		result = strings.Join(errorslice, ", ")
		return "", errors.New(result)
	}

checks:
	if session.State != AUTHORIZATION {
		errorslice = append(errorslice, "wrong session state")
	}
	if session.Authenticated {
		errorslice = append(errorslice, "already authenticated")
	}
	if session.Username != "" {
		errorslice = append(errorslice, "user already set")
	}
	if len(c.Arguments) == 1 {
		errorslice = append(errorslice, "username can't be empty")
	}
	if len(c.Arguments) == 2 {
		errorslice = append(errorslice, "digest can't be empty")
	}
	if len(c.Arguments) > 3 {
		errorslice = append(errorslice, "too many arguments")
	}

	if len(errorslice) != 0 {
		goto returnerror
	}

	log.Println("POP3:", "APOP command issued by", session.RemoteEP, "with", c.Arguments[1])

	user, erra := mailboxes.GetUser(c.Arguments[1])
	if erra != nil {
		errorslice = append(errorslice, "no such user")
		goto returnerror
	}

	hash := md5.New()
	hash.Write([]byte(session.Shared + user.Arguments[3]))
	digest := hex.EncodeToString(hash.Sum(nil))

	if c.Arguments[2] == digest {
		lockerr := locker.Lock(mailboxes.GetMailbox(session.Username))
		if lockerr != nil {
			errorslice = append(errorslice, "[IN-USE] maildrop "+lockerr.Error())
			goto returnerror
		}

		session.Username = c.Arguments[1]
		session.Password = c.Arguments[1]
		session.Authenticated = true
		session.State = TRANSACTION
		count, octets := mailboxes.Stat(session.Username, false)
		result = session.Username + "'s maildrop has " + strconv.Itoa(count) + " messages (" + strconv.Itoa(octets) + " octets)"
	} else {
		errorslice = append(errorslice, "invalid user/digest combination")
		goto returnerror
	}

	return result, nil
}
