package list

import (
	"errors"
	. "github.com/trapped/gomaild/parsers/textual"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/stat"
	. "github.com/trapped/gomaild/processors/pop3/session"
	"log"
	"strconv"
)

func Process(session *Session, c Command) (string, error) {
	if session.Username == "" || session.Password == "" || session.Authenticated == false || session.State != TRANSACTION {
		return "", errors.New("wrong state")
	} else {
		log.Println("POP3: Attempt to use LIST as", session.Username)
		if len(c.Arguments) == 1 {
			count, moctets := stat.Stat(session)
			result := strconv.Itoa(count) + " messages (" + strconv.Itoa(moctets) + " octets)\r\n"
			messages := stat.Index(session)
			for i := 0; i < len(messages); i++ {
				result += strconv.Itoa(messages[i].ID) + " " + strconv.Itoa(int(messages[i].File.Size())) + "\r\n"
			}
			result += "."
			return result, nil
		} else {
			messages := stat.Index(session)
			for _, v := range messages {
				if strconv.Itoa(v.ID) == c.Arguments[1] {
					return strconv.Itoa(v.ID) + " " + strconv.Itoa(int(v.File.Size())), nil
				}
			}
			return "", errors.New("no such message, " + strconv.Itoa(len(messages)) + " messages in maildrop")
		}
	}
}
