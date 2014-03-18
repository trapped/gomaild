package pop3client

import (
	"bufio"
	"github.com/trapped/gomaild/config"
	"github.com/trapped/gomaild/parsers/textual"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor"
	"github.com/trapped/gomaild/processors/pop3/locker"
	"github.com/trapped/gomaild/processors/pop3/session"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

type Client struct {
	Parent   *net.Listener
	Conn     net.Conn
	Start    time.Time
	End      time.Time
	KeepOpen bool
}

func MakeClient(parent *net.Listener, conn net.Conn) *Client {
	return &Client{
		Parent:   parent,
		Conn:     conn,
		Start:    time.Now(),
		KeepOpen: true,
	}
}

func (c *Client) RemoteEP() string {
	return c.Conn.RemoteAddr().String()
}

func (c *Client) LocalEP() string {
	return c.Conn.LocalAddr().String()
}

func (c *Client) Process() {
	defer c.Conn.Close()
	bufin := bufio.NewReader(c.Conn)
	processor := cmdprocessor.Processor{
		Session: &session.Session{},
	}
	greeting := "+OK"
	if config.Settings["pop3"] != nil && len(config.Settings["pop3"]["greeting"]) >= 1 && strings.TrimSpace(config.Settings["pop3"]["greeting"][0].(textual.Command).Arguments[1]) != "" {
		greeting += " " + config.Settings["pop3"]["greeting"][0].(textual.Command).Arguments[1]
	}
	_, errX := c.Conn.Write([]byte(greeting + "\r\n"))
	if errX != nil {
		log.Println(errX)
		return
	}
	for c.KeepOpen {
		if processor.Session.Quitted {
			break
		}
		line, err := bufin.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}
		_, err0 := c.Conn.Write([]byte(processor.Process(line)))
		if err0 != nil {
			log.Println(err0)
			return
		}
	}
	locker.Unlock(path.Dir(os.Args[0]) + "/mailboxes/" + processor.Session.Username)
	c.End = time.Now()
}
