package mailboxes

import (
	"errors"
	"github.com/trapped/gomaild/config"
	. "github.com/trapped/gomaild/parsers/textual"
	"os"
	"path"
	"path/filepath"
)

func GetUser(s string) (Statement, error) {
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
	return path.Dir(os.Args[0]) + "/mailboxes/" + user
}

func Stat(user string, showdeleted bool) (int, int) {
	count, octets := 0, 0

	walkFn := func(p string, info os.FileInfo, e error) error {
		if !info.IsDir() && e == nil {
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
