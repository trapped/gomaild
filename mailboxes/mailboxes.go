//Package mailboxes provides static utility functions to manage mailboxes.
package mailboxes

import (
	"errors"
	"github.com/trapped/gomaild/config"
	"log"
	"os"
	"path"
	"path/filepath"
)

//Checks the existence of a username; if it exists, returns its password, otherwise returns an error.
func GetUser(s string) (string, error) {
	log.Println("package/mailboxes:", "Queried existence of user <"+s+">")
	if value, exists := config.Configuration.Accounts[s]; exists {
		return value, nil
	}
	return "", errors.New("no such user")
}

//Blindly generates the mailbox filepath for a username. Check the existence of the username with GetUser first.
func GetMailbox(user string) string {
	log.Println("package/mailboxes:", "Queried blind path of mailbox for user <"+user+">")
	return path.Dir(os.Args[0]) + "/mailboxes/" + user
}

//Blindly checks the number of files in a username's mailbox, their complessive size. It's possible to decide whether to include emails in the "deleted" folder.
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

//Checks the existence of a username's mailbox, and, if it doesn't exist, creates it.
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

//Creates a mailbox for the given username.
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
