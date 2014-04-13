package pop3

import (
	. "github.com/trapped/gomaild/processors/pop3/client"
	"log"
	"net"
	"os"
	"strconv"
)

type POP3 struct {
	Port int  //Port to listen at
	Keep bool //Whether to keep accepting clients (doesn't prevent active clients from continuing their current sessions)
}

func (p *POP3) Listen() {
	log.Println("POP3: Starting POP3 server")
	if p.Keep == false {
		p.Keep = true
	}
	//Set the default port if the user hasn't specified one
	if p.Port == 0 {
		p.Port = 110
	}
	//Start listening on all interfaces and the specified port
	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(p.Port))
	if err != nil {
		log.Println("POP3:", err)
		os.Exit(1)
	}
	for p.Keep {
		client, err := listener.Accept()
		if err != nil {
			log.Println("POP3:", err)
			continue //Ignore errors, we must remain available as much as possible
		}
		cliobj := MakeClient(&listener, client)
		log.Println("POP3: Accepting client", cliobj.RemoteEP())
		go cliobj.Process()
	}
}
