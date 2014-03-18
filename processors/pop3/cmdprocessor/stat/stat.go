package stat

import (
	"errors"
	. "github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/pop3/session"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func Process(session *Session, c Command) (string, error) {
	if !session.Authenticated || session.Username == "" || session.Password == "" || session.State != TRANSACTION {
		return "", errors.New("client not authenticated")
	}
	log.Println("POP3: Attempt to STAT by", session.Username)
	count, octets := Stat(session)
	return strconv.Itoa(count) + " " + strconv.Itoa(octets), nil
}

type Messages []Message

type Message struct {
	ID   int
	File os.FileInfo
	Path string
}

func (m Messages) Len() int {
	return len(m)
}

func (m Messages) Less(i, j int) bool {
	return m[i].File.ModTime().Before(m[j].File.ModTime())
}

func (m Messages) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func Index(session *Session) []Message {
	files := make([]Message, 0)
	walkfn := func(p string, info os.FileInfo, e error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".eml") && filepath.Base(filepath.Dir(p)) != "deleted" && e == nil {
			files = append(files, Message{File: info, Path: p})
		}

		return nil
	}
	filepath.Walk(path.Dir(os.Args[0])+"/mailboxes/"+session.Username, walkfn)
	sort.Sort(Messages(files))
	for i, _ := range files {
		files[i].ID = i + 1
	}
	return files
}

func Stat(session *Session) (int, int) {
	messages := Index(session)
	wholesize := 0
	for i := 0; i < len(messages); i++ {
		wholesize += int(messages[i].File.Size())
	}
	return len(messages), wholesize
}
