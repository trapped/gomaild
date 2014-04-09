package mailboxes

import (
	"errors"
	"github.com/trapped/gomaild/config"
	. "github.com/trapped/gomaild/parsers/textual"
	"log"
	"os"
	"path"
	"path/filepath"
)

func GetUser(s string) (Statement, error) {
	log.Println("package/mailboxes:", "Queried existence of user <"+s+">")
	if config.Settings["gomaild"] != nil {
		for _, v := range config.Settings["gomaild"]["user"] {
			z := v.(Statement)
			if z.Arguments[1] == s {
				return z, nil
			}
		}
	}
	if config.Settings["pop3"] != nil {
		for _, v := range config.Settings["pop3"]["user"] {
			z := v.(Statement)
			if z.Arguments[1] == s {
				return z, nil
			}
		}
	}
	return Statement{}, errors.New("no such user")
}

func GetMailbox(user string) string {
	log.Println("package/mailboxes:", "Queried blind path of mailbox for user <"+user+">")
	return path.Dir(os.Args[0]) + "/mailboxes/" + user
}

func Stat(user string, showdeleted bool) (int, int) {
	log.Println("package/mailboxes:", "Queried stat[", showdeleted, "] of mailbox for user <"+user+">")
	count, octets := 0, 0

	walkFn := func(p string, info os.FileInfo, e error) error {
		if e == nil && !info.IsDir() {
			if path.Base(path.Dir(p)) == "deleted" && !showdeleted {
				return nil
			}

			count++
			octets += int(info.Size())
		}

		return nil
	}

	filepath.Walk(GetMailbox(user), walkFn)

	return count, octets
}

func CreateIfNull(user string) {
	log.Println("package/mailboxes:", "Checking existence of mailbox <"+user+">")
	if _, err := os.Stat(GetMailbox(user)); err != nil {
		if !os.IsNotExist(err) {
			log.Println("package/mailboxes:", "Error checking the existence of the <"+user+"> mailbox:", err)
		}
		errc := Create(user)
		if errc != nil {
			log.Println("package/mailboxes:", "Failed creating mailbox for user <"+user+">:", err)
		}
	}
}

func Create(user string) error {
	log.Println("package/mailboxes:", "Creating mailbox for user <"+user+">")
	err := os.Mkdir(GetMailbox(user), 0777)
	if err != nil {
		return err
	}
	err = os.Mkdir(GetMailbox(user)+"/unread", 0777)
	if err != nil {
		return err
	}
	err = os.Mkdir(GetMailbox(user)+"/read", 0777)
	if err != nil {
		return err
	}
	err = os.Mkdir(GetMailbox(user)+"/deleted", 0777)
	if err != nil {
		return err
	}
	return nil
}
