package pass

import (
	"errors"
	. "github.com/trapped/gomaild/parsers/textual"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/stat"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/user"
	"github.com/trapped/gomaild/processors/pop3/locker"
	. "github.com/trapped/gomaild/processors/pop3/session"
	"log"
	"os"
	"path"
	"strconv"
)

func Process(session *Session, c Command) (string, error) {
	if !session.Authenticated && session.Username != "" && session.Password == "" && session.State == AUTHORIZATION {
		log.Println("POP3: Attempt to use PASS", c.Arguments[1], "with USER", session.Username)
		cmd, err := user.GetUser(session.Username)
		if err != nil || cmd.Arguments[3] != c.Arguments[1] {
			session.Username = ""
			return "", errors.New("incorrect username/password combination")
		} else {
			if errl := locker.Lock(path.Dir(os.Args[0]) + "/mailboxes/" + session.Username); errl != nil {
				session.Username = ""
				return "", errors.New("maildrop " + errl.Error())
			}
			session.Password = c.Arguments[1]
			session.Authenticated = true
			session.State = TRANSACTION
			count, octets := stat.Stat(session)
			return session.Username + "'s maildrop has " + strconv.Itoa(count) + " messages (" + strconv.Itoa(octets) + " octets)", nil
		}
	}
	session.Username = ""
	session.Password = ""
	return "", errors.New("wrong state")
}
