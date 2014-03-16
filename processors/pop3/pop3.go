package pop3

import (
	. "github.com/trapped/gomaild/processors/pop3/client"
	"log"
	"net"
	"os"
	"strconv"
)

type POP3 struct {
	//Port to listen at
	Port int
	//Whether to keep accepting clients (doesn't prevent active clients from continuing their current sessions)
	Keep bool
}

func (p *POP3) Listen() {
	if p.Keep == false {
		p.Keep = true
	}
	if p.Port == 0 {
		p.Port = 110
	}
	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(p.Port))
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	for p.Keep {
		client, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		cliobj := MakeClient(&listener, client)
		go cliobj.Process()
	}
}
