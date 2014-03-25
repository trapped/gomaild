package uidl

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	. "github.com/trapped/gomaild/parsers/textual"
	"github.com/trapped/gomaild/processors/pop3/message"
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
		result = strings.Join(errorslice, ", ")
		return "", errors.New(result)
	}

checks:
	if session.State != TRANSACTION {
		errorslice = append(errorslice, "wrong session state")
	}
	if !session.Authenticated {
		errorslice = append(errorslice, "not authenticated")
	}
	if session.Username == "" {
		errorslice = append(errorslice, "user can't be empty")
	}
	if len(c.Arguments) > 2 {
		errorslice = append(errorslice, "too many arguments")
	}

	if len(errorslice) != 0 {
		goto returnerror
	}

	log.Println("POP3:", "UIDL command issued by", session.RemoteEP, "with", session.Username)

	messages := message.Index(session)
	if len(c.Arguments) == 1 {
		result = "\r\n"
		for i := 0; i < len(messages); i++ {
			result += strconv.Itoa(messages[i].ID) + " "
			messageuid := ""
			for f := 0; f < 5; f++ {
				part := messages[i].Path + strconv.Itoa(messages[i].ID) + messages[i].File.ModTime().String()
				hash := md5.New()
				hash.Write([]byte(part))
				messageuid += hex.EncodeToString(hash.Sum(nil))
			}
			result += messageuid[:69]
			result += "\r\n"
		}
		result += "."
	} else {
		for _, v := range messages {
			if strconv.Itoa(v.ID) == c.Arguments[1] {
				result = strconv.Itoa(v.ID) + " "
				messageuid := ""
				for f := 0; f < 5; f++ {
					part := v.Path + strconv.Itoa(v.ID) + v.File.ModTime().String()
					hash := md5.New()
					hash.Write([]byte(part + messageuid))
					hex.EncodeToString(hash.Sum(nil))
					messageuid += hex.EncodeToString(hash.Sum(nil))
				}
				result += messageuid[:69]
				break
			}
		}
		if result == "" {
			errorslice = append(errorslice, "no such message; "+strconv.Itoa(len(messages))+" messages in maildrop")
			goto returnerror
		}
	}

	return result, nil
}
