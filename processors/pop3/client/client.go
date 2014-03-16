package pop3client

import (
	"bufio"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor"
	"log"
	"net"
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
	processor := cmdprocessor.Processor{CommandLock: false}
	for c.KeepOpen {
		line, err := bufin.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}
		_, err0 := c.Conn.Write([]byte(processor.Process(line)))
		if err0 != nil {
			log.Println(err)
			return
		}
	}
	c.End = time.Now()
}
