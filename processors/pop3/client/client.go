//Package client defines a set of structs and methods to handle a POP3 client.
package client

import (
	"bufio"
	"github.com/trapped/gomaild/config"
	"github.com/trapped/gomaild/locker"
	"github.com/trapped/gomaild/mailboxes"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor"
	"github.com/trapped/gomaild/processors/pop3/session"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

//Stores objects and data necessary to handle a POP3 client.
type Client struct {
	Parent       *net.Listener //The client's parent listener.
	Conn         net.Conn      //The client's network connection.
	Start        time.Time     //The connection start time.
	End          time.Time     //The connection end time.
	KeepOpen     bool          //Whether to keep Process() looping or not.
	TimeoutTimer *time.Timer   //Timer to check whether the connection times out.
}

//Creates a Client struct ready to be used, given a parent network listener and a network connection.
func MakeClient(parent *net.Listener, conn net.Conn) *Client {
	return &Client{
		Parent:       parent,
		Conn:         conn,
		Start:        time.Now(),
		KeepOpen:     true,
		TimeoutTimer: time.NewTimer(time.Duration(config.Configuration.POP3.Timeout) * time.Second),
	}
}

//Returns a string representing a client's connection remote endpoint, complete of IP address and port.
func (c *Client) RemoteEP() string {
	return c.Conn.RemoteAddr().String()
}

//Returns a string representing a client's connection local endpoint, complete of IP address and port.
func (c *Client) LocalEP() string {
	return c.Conn.LocalAddr().String()
}

//Sends a line of text, usually obtained from the execution of a POP3 command.
//Already appends the CRLF (0x0a 0x0d) termination octet pair.
func (c *Client) Send(s string) error {
	_, err := c.Conn.Write([]byte(s + "\r\n"))
	return err
}

//Loops and processes the client's commands, until something changes its Client's KeepOpen property to false, the client QUITs the session, or it encounters an error trying to read/write on the connection.
func (c *Client) Process() {
	//Close the network connection at the end of the function if it's not closed already.
	defer c.Conn.Close()

	//Set up a buffered-IO binary reader.
	bufin := bufio.NewReader(c.Conn)
	//Set up a command processor.
	processor := cmdprocessor.Processor{
		Session: &session.Session{RemoteEP: c.RemoteEP()},
	}

	//Set the POP3 session unique shared
	processor.Session.Shared = "<" + strconv.Itoa(os.Getpid()) + "." + strconv.Itoa(time.Now().Nanosecond()) + "@" + config.Configuration.ServerName + ">"

	//Send the POP3 session-start greeting eventually set in the "pop3.conf" configuration file.
	err1 := c.Send("+OK " + config.Configuration.POP3.StartGreeting + " " + processor.Session.Shared)
	//If an error occurs, log it and finalize the connection.
	if err1 != nil {
		log.Println("POP3:", err1)
		return
	}

	//Start the goroutine to check for timeout
	go c.Timeout()

	//Stop looping if the KeepOpen property of the client becomes false.
	for c.KeepOpen {
		//Stop looping if the client quits its POP3 session.
		if processor.Session.Quitted {
			break
		}

		//Read a line from the network stream.
		//The POP3 protocol defines the command termination octet pair as CRLF (0x0d 0x0a); however, we simply read until a LF occurs (the termination octet pair is stripped away later).
		line, err2a := bufin.ReadString('\n')
		//If an error occurs, log it and finalize the connection.
		if err2a != nil {
			log.Println(err2a)
			break
		}

		//Reset the timeout timer
		c.TimeoutTimer.Reset(time.Duration(config.Configuration.POP3.Timeout) * time.Second)

		//Process the last command received from the client using cmdprocessor.Process() and send the result to the client.
		err2b := c.Send(processor.Process(line))
		//If an error occurs, log it and finalize the connection.
		if err2b != nil {
			log.Println(err2b)
			break
		}
	}

	//If the client opened a mailbox, unlock it using the locker.
	if processor.Session.Authenticated {
		locker.Unlock(mailboxes.GetMailbox(processor.Session.Username))
	}

	//Set the Client's connection end time to time.Now().
	c.End = time.Now()

	//Log the connection ending, with the remote endpoint.
	log.Println("POP3: Disconnecting", c.RemoteEP())
}

//Checks its client handler for POP3 timeout.
func (c *Client) Timeout() {
	for c.KeepOpen {
		select {
		case x := <-c.TimeoutTimer.C:
			log.Println("POP3:", "Timeout for", c.RemoteEP()+":", x)
			c.End = time.Now()
			c.Send("-ERR " + config.Configuration.POP3.TimeoutMessage)
			c.KeepOpen = false
			err := c.Conn.Close()
			if err != nil {
				log.Println("POP3:", "Error closing the connection for", c.RemoteEP()+":", err)
			}
			return
		}
	}
}
