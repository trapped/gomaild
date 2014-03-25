//Package client defines a set of structs and methods to handle a SMTP client.
package client

import (
	"bufio"
	"github.com/trapped/gomaild/locker"
	"github.com/trapped/gomaild/processors/smtp/cmdprocessor"
	"github.com/trapped/gomaild/processors/smtp/sentences"
	"github.com/trapped/gomaild/processors/smtp/session"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"time"
)

//Client is the structure used to store some useful parts of the SMTP clients.
type Client struct {
	Parent   *net.Listener //The client's parent listener.
	Conn     net.Conn      //The client's network connection.
	Start    time.Time     //The connection start time.
	End      time.Time     //The connection end time.
	KeepOpen bool          //Whether to keep Process() looping or not.
}

//MakeClient creates a client, given a parent network listener and a network connection.
func MakeClient(parent *net.Listener, conn net.Conn) *Client {
	return &Client{
		Parent:   parent,
		Conn:     conn,
		Start:    time.Now(),
		KeepOpen: true,
	}
}

//RemoteEP returns a string representing a client's connection remote endpoint, complete of IP address and port.
func (c *Client) RemoteEP() string {
	return c.Conn.RemoteAddr().String()
}

//LocalEP returns a string representing a client's connection local endpoint, complete of IP address and port.
func (c *Client) LocalEP() string {
	return c.Conn.LocalAddr().String()
}

//Send sends a line of text, usually obtained from the execution of a SMTP command.
//Send already appends the CRLF (0x0a 0x0d) termination octet pair.
func (c *Client) Send(s string) error {
	_, err := c.Conn.Write([]byte(s + "\r\n"))
	return err
}

//Process loops and processes the client's commands.
//Process loops until: something changes its Client's KeepOpen property to false, the client QUITs the session, it encounters an error trying to read/write on the connection.
func (c *Client) Process() {
	//Close the network connection at the end of the function if it's not closed already.
	defer c.Conn.Close()

	//Set up a buffered-IO binary reader.
	bufin := bufio.NewReader(c.Conn)
	//Set up a command processor.
	processor := cmdprocessor.Processor{
		Session: &session.Session{RemoteEP: c.RemoteEP()},
	}

	//Set the SMTP session unique shared
	processor.Session.Shared = "<" + strconv.Itoa(os.Getpid()) + "." + strconv.Itoa(time.Now().Nanosecond()) + ">"

	//Send the SMTP session-start greeting eventually set in the "smtp.conf" configuration file and the shared.
	err1 := c.Send("220 " + sentences.StartGreeting() + " " + processor.Session.Shared)
	//If an error occurs, log it and finalize the connection.
	if err1 != nil {
		log.Println("SMTP:", err1)
		return
	}

	//Stop looping if the KeepOpen property of the client becomes false.
	for c.KeepOpen {
		//Stop looping if the client quits its SMTP session.
		if processor.Session.Quitted {
			break
		}

		//Read a line from the network stream.
		//The SMTP protocol defines the command termination octet pair as CRLF (0x0d 0x0a); however, we simply read until a LF occurs (the termination octet pair is stripped away later).
		line, err2a := bufin.ReadString('\n')
		//If an error occurs, log it and finalize the connection.
		if err2a != nil {
			log.Println(err2a)
			break
		}

		//Process the last command received from the client using cmdprocessor.Process() and send the result to the client.
		err2b := c.Send(processor.Process(line))
		//If an error occurs, log it and finalize the connection.
		if err2b != nil {
			log.Println(err2b)
			break
		}
	}

	//If the client opened a mailbox, unlock it using the locker.
	//if processor.Session.Authenticated {
	//	locker.Unlock(path.Dir(os.Args[0]) + "/mailboxes/" + processor.Session.Username)
	//}

	//Set the Client's connection end time to time.Now().
	c.End = time.Now()

	//Log the connection ending, with the remote endpoint.
	log.Println("SMTP: Disconnecting", c.RemoteEP())
}
