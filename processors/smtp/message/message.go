package message

import (
	"github.com/trapped/gomaild/locker"
	"github.com/trapped/gomaild/mailboxes"
	. "github.com/trapped/gomaild/processors/smtp/session"
	rfc "github.com/trapped/rfc2822"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
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

func Store(session *Session, m Message) error {
	log.Println("SMTP/message:", "Storing message for", m.Recipients, "by", m.Sender)
	msg, err := rfc.ReadString(m.Text)
	if err != nil {
		log.Println("SMTP/message:", "Error parsing message:", err)
	}
	msg.AddMultiHeader("Received", []string{"from " + m.RemoteDomain + "(" + session.RemoteEP + ")" + ";", time.Now().Format(time.RFC1123Z)})
	if m.Sender != "" {
		msg.AddHeader("Return-Path", "<"+m.Sender+">")
	}
	if len(m.Recipients) != 0 {
		msg.AddHeader("X-Recipients", func(a []string) string {
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

	newfull := msg.Text()

	for _, v := range m.Recipients {
		mailboxes.CreateIfNull(v)
		locker.MLock(mailboxes.GetMailbox(v))
		ustatc, _ := mailboxes.Stat(v, true)
		err := ioutil.WriteFile(mailboxes.GetMailbox(v)+"/unread/"+strconv.Itoa(ustatc)+".eml", []byte(newfull), 0777)
		if err != nil {
			return err
		}
		locker.MUnlock(mailboxes.GetMailbox(v))
	}

	return nil
}
