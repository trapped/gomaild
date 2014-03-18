package user

import (
	"errors"
	"github.com/trapped/gomaild/config"
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/pop3/session"
	"log"
	"strconv"
)

func Process(session *Session, c Command) (string, error) {
	if !session.Authenticated && session.Username == "" && session.Password == "" && session.State == AUTHORIZATION {
		log.Println("POP3: Attempt to use USER", c.Arguments[1])
		success := false
		if c.Arguments[1] != "" {
			if config.Settings["pop3"]["secure_user"] != nil || len(config.Settings["pop3"]["secure_user"]) >= 1 {
				if su, _ := strconv.ParseBool(config.Settings["pop3"]["secure_user"][0].(Command).Arguments[1]); su {
					success = true
				}
			} else {
				_, err := GetUser(c.Arguments[1])
				if err != nil {
					success = false
				} else {
					success = true
				}
			}
			if success {
				session.Username = c.Arguments[1]
				return "user may exist", nil
			} else {
				return "", errors.New("no such user")
			}
		} else {
			return "", errors.New("username can't be empty")
		}
	}
	session.Username = ""
	return "", errors.New("wrong state")
}

func GetUser(s string) (Command, error) {
	if config.Settings["gomaild"] != nil {
		for _, v := range config.Settings["gomaild"]["user"] {
			z := v.(Command)
			if z.Arguments[1] == s {
				return z, nil
			}
		}
	}
	if config.Settings["pop3"] != nil {
		for _, v := range config.Settings["pop3"]["user"] {
			z := v.(Command)
			if z.Arguments[1] == s {
				return z, nil
			}
		}
	}
	return Command{}, errors.New("no such user")
}
