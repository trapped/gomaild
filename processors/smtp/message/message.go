//Package message provides SMTP utility functions for MIME messages.
package message

import (
	"fmt"
	"github.com/trapped/gomaild/config"
	"github.com/trapped/gomaild/locker"
	"github.com/trapped/gomaild/mailboxes"
	. "github.com/trapped/gomaild/processors/smtp/session"
	rfc "github.com/trapped/rfc2822"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

//type Messages []Message

//A single message.
type Message struct {
	ID           int         //The non-unique ID of the message, if picked up from disk
	File         os.FileInfo //The FileInfo of the message, if picked up from disk
	Path         string      //The filepath of the message, if picked up from disk
	RemoteDomain string      //The identity used with EHLO or HELO
	Sender       string      //Return-Path set with the MAIL fommand
	Recipients   []string    //Recipients set with various RCPT commands
	Text         string      //Data received with the DATA command
}

//Stores a Message on disk, adding useful headers ("Received", "Return-Path").
func Store(session *Session, m Message) error {
	log.Println("SMTP/message:", "Storing message for", m.Recipients, "by", m.Sender)

	//Parse message
	msg, err := rfc.ReadString(m.Text)
	if err != nil {
		log.Println("SMTP/message:", "Error parsing message:", err)
	}

	//Lookup the client's text-form address
	remdom, err := net.LookupAddr(strings.TrimPrefix(strings.Split(session.RemoteEP, "]")[0], "["))
	if err != nil {
		log.Println("SMTP/message:", "Failed looking up address for", session.RemoteEP)
		remdom = append(remdom, session.RemoteEP)
	}

	//Add the "Received" header
	recvfrom := fmt.Sprintf("from %s (EHLO %s) (%s)", remdom[0], session.Identity, session.RemoteEP)
	if session.Authenticated {
		recvfrom += fmt.Sprintf("(AUTH %s %s)", strings.ToUpper(session.AuthMode), session.Username)
	}
	recvby := fmt.Sprintf("by %s using gomaild", config.Configuration.ServerName)
	recvfor := fmt.Sprintf("for <%s>;", strings.Join(m.Recipients, ", "))
	recvdate := time.Now().Format(time.RFC1123Z)
	msg.AddMultiHeader("Received", []string{recvfrom, recvby, recvfor, recvdate})
	//Add "Return-Path" header
	if m.Sender != "" {
		msg.AddHeader("Return-Path", "<"+m.Sender+">")
	}

	//Merge headers and body into a single, full message
	newfull := msg.Text()

	//Save the message in each of the recipients' mailboxes, in the "unread" folder
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
