package message

import (
	"github.com/trapped/gomaild/locker"
	"github.com/trapped/gomaild/mailboxes"
	. "github.com/trapped/gomaild/processors/smtp/session"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Messages []Message

type Message struct {
	ID           int
	File         os.FileInfo
	Path         string
	RemoteDomain string
	Sender       string
	Recipients   []string
	Text         string
}

func Headers(m Message) (string, error) {
	file := []byte{}
	if m.Text == "" {
		afile, err := ioutil.ReadFile(m.Path)
		if err != nil {
			return "", err
		}
		file = afile
	} else {
		file = []byte(m.Text)
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

func HeadersLimit(m Message) (int, error) {
	file := []byte{}
	if m.Text == "" {
		afile, err := ioutil.ReadFile(m.Path)
		if err != nil {
			return 0, err
		}
		file = afile
	} else {
		file = []byte(m.Text)
	}
	lines := strings.Split(string(file), "\r\n")
	for i, v := range lines {
		if v == "" {
			return i, nil
		}
	}
	return 0, nil
}

func Body(m Message) (string, error) {
	file := []byte{}
	if m.Text == "" {
		afile, err := ioutil.ReadFile(m.Path)
		if err != nil {
			return "", err
		}
		file = afile
	} else {
		file = []byte(m.Text)
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

func Store(session *Session, m Message) error {
	hders, err := Headers(m)
	if err != nil {
		return err
	}
	headers := strings.Split(hders, "\r\n")
	if m.RemoteDomain != "" {
		headers = append(headers, "X-Received-From: "+m.RemoteDomain)
	}
	if m.Sender != "" {
		headers = append(headers, "X-Sender: <"+m.Sender+">")
	}
	if len(m.Recipients) != 0 {
		headers = append(headers, "X-Recipient: "+func(a []string) string {
			result := ""
			for _, v := range a {
				if result != "" {
					result += ","
				}
				result += "<" + v + ">"
			}
			return result
		}(m.Recipients))
	}
	mheaders := strings.Join(headers, "\r\n")
	mbody, err := Body(m)
	if err != nil {
		return err
	}
	newfull := mheaders + "\r\n" + mbody

	for _, v := range m.Recipients {
		locker.MLock(mailboxes.GetMailbox(v))
		ustatc, _ := mailboxes.Stat(v, true)
		ioutil.WriteFile(mailboxes.GetMailbox(v)+"/unread/"+strconv.Itoa(ustatc)+".eml", []byte(newfull), 0777)
		locker.MUnlock(mailboxes.GetMailbox(v))
	}

	return nil
}
