//Package client defines a set of structs and methods to handle a SMTP client.
package client

import (
	"bufio"
	"github.com/trapped/gomaild/cipher"
	"github.com/trapped/gomaild/config"
	"github.com/trapped/gomaild/processors/smtp/cmdprocessor"
	. "github.com/trapped/gomaild/processors/smtp/session"
	"log"
	"net"
	"time"
)

//Stores objects and data necessary to handle SMTP clients.
type Client struct {
	Parent       *net.Listener //The client's parent listener.
	Conn         net.Conn      //The client's network connection.
	Start        time.Time     //The connection start time.
	End          time.Time     //The connection end time.
	KeepOpen     bool          //Whether to keep Process() looping or not.
	TimeoutTimer *time.Timer   //Timer to check whether the connection times out.
}

//Creates a Client ready to be used, given a parent network listener and a network connection.
func MakeClient(parent *net.Listener, conn net.Conn) *Client {
	return &Client{
		Parent:       parent,
		Conn:         conn,
		Start:        time.Now(),
		KeepOpen:     true,
		TimeoutTimer: time.NewTimer(time.Duration(config.Configuration.SMTP.Timeout) * time.Second),
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

//Sends a line of text, usually obtained from the execution of a SMTP command. Already appends the CRLF (0x0a 0x0d) termination octet pair.
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
		Session: &Session{RemoteEP: c.RemoteEP()},
	}

	//Send the SMTP session-start greeting eventually set in the "smtp.conf" configuration file and the shared.
	err1 := c.Send("220 " + config.Configuration.SMTP.StartGreeting /* + " " + processor.Session.Shared*/)
	//If an error occurs, log it and finalize the connection.
	if err1 != nil {
		log.Println("SMTP:", err1)
		return
	}

	//Start the goroutine to check for timeout
	go c.Timeout()

	for c.KeepOpen {
		//Stop looping if the client quits its SMTP session.
		if processor.Session.Quitted {
			c.KeepOpen = false
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

		//Reset the timeout timer
		c.TimeoutTimer.Reset(time.Duration(config.Configuration.SMTP.Timeout) * time.Second)

		//Process the last command received from the client using cmdprocessor.Process() and send the result to the client. If the processor is waiting for a multiline message, just wait until it exits the COMPOSITION state.
		oldstate := processor.Session.State
		oldtls := processor.Session.InTLS
		result := processor.Process(line)
		if processor.Session.State == COMPOSITION && oldstate == COMPOSITION {
			continue
		}
		err2b := c.Send(result)
		//If an error occurs, log it and finalize the connection.
		if err2b != nil {
			log.Println(err2b)
			break
		}
		if processor.Session.InTLS && !oldtls {
			c.Conn = cipher.TLSTransmuteConn(c.Conn)
			bufin = bufio.NewReader(c.Conn)
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

//Checks its client for SMTP timeout.
func (c *Client) Timeout() {
	for c.KeepOpen {
		select {
		case x := <-c.TimeoutTimer.C:
			log.Println("SMTP:", "Timeout for", c.RemoteEP()+":", x)
			c.End = time.Now()
			c.Send("421 " + config.Configuration.SMTP.TimeoutMessage)
			c.KeepOpen = false
			err := c.Conn.Close()
			if err != nil {
				log.Println("SMTP:", "Error closing the connection for", c.RemoteEP()+":", err)
			}
			return
		}
	}
}
