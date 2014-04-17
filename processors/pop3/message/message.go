//Package message provides POP3 utility functions for MIME messages.
package message

import (
	"github.com/trapped/gomaild/mailboxes"
	. "github.com/trapped/gomaild/processors/pop3/session"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//Moves the message from a folder to another inside a user's mailbox.
func MoveMessage(session *Session, m Message, destfolder string) error {
	d := mailboxes.GetMailbox(session.Username) + "/" + destfolder + "/" + m.File.Name()
	err := os.Rename(m.Path, d)
	if err != nil {
		return err
	}
	return nil
}

//Deletes permanently the message.
func DeleteMessage(m Message) error {
	err := os.Remove(m.Path)
	if err != nil {
		return err
	}
	return nil
}

//Type necessary for sorting.
type Messages []Message

//A single message.
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

//Returns a sorted array of messages inside a user's mailbox.
func Index(session *Session) []Message {
	files := make([]Message, 0)
	walkFn := func(p string, info os.FileInfo, e error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".eml") && filepath.Base(filepath.Dir(p)) != "deleted" && e == nil {
			files = append(files, Message{File: info, Path: p})
		}

		return nil
	}
	filepath.Walk(mailboxes.GetMailbox(session.Username), walkFn)
	sort.Sort(Messages(files))
	for i, _ := range files {
		files[i].ID = i + 1
	}
	return files
}

func MessagesContain(i []interface{}, id int) bool {
	for q, _ := range i {
		if i[q].(Message).ID == id {
			return true
		}
	}
	return false
}

//Returns the headers of the message.
func Headers(m Message) (string, error) {
	file, err := ioutil.ReadFile(m.Path)
	if err != nil {
		return "", err
	}
	pos, err := HeadersLimit(m)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(file), "\r\n")
	headers := []string{""}
	if pos > 0 {
		headers = lines[:pos]
	}
	return strings.Join(headers, "\r\n"), nil
}

//Returns the line number of the separator between the headers and the body of the message.
func HeadersLimit(m Message) (int, error) {
	file, err := ioutil.ReadFile(m.Path)
	if err != nil {
		return 0, err
	}
	lines := strings.Split(string(file), "\r\n")
	for i, v := range lines {
		if v == "" {
			return i, nil
		}
	}
	return 0, nil
}

//Returns the body of the message.
func Body(m Message) (string, error) {
	file, err := ioutil.ReadFile(m.Path)
	if err != nil {
		return "", err
	}
	pos, err := HeadersLimit(m)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(file), "\r\n")
	body := []string{""}
	if pos <= len(lines)-1 {
		body = lines[pos+1:]
	}
	return strings.Join(body, "\r\n"), nil
}
