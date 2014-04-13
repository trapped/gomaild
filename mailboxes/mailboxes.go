package mailboxes

import (
	"errors"
	"github.com/trapped/gomaild/config"
	"log"
	"os"
	"path"
	"path/filepath"
)

func GetUser(s string) (string, error) {
	log.Println("package/mailboxes:", "Queried existence of user <"+s+">")
	if value, exists := config.Configuration.Accounts[s]; exists {
		return value, nil
	}
	return "", errors.New("no such user")
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
