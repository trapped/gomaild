package smtp

import (
	. "github.com/trapped/gomaild/processors/smtp/client"
	"log"
	"net"
	"os"
	"strconv"
)

type SMTP struct {
	//Port to listen at
	Port int
	//Whether to keep accepting clients (doesn't prevent active clients from continuing their current sessions)
	Keep bool
}

func (p *SMTP) Listen() {
	log.Println("SMTP: Starting SMTP server")
	if p.Keep == false {
		p.Keep = true
	}
	if p.Port == 0 {
		p.Port = 25
	}

	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(p.Port))
	if err != nil {
		log.Println("SMTP:", err)
		os.Exit(1)
	}
	for p.Keep {
		client, err := listener.Accept()
		if err != nil {
			log.Println("SMTP:", err)
			continue
		}
		cliobj := MakeClient(&listener, client)
		log.Println("SMTP: Accepting client", cliobj.RemoteEP())
		go cliobj.Process()
	}
}
