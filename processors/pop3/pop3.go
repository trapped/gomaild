package pop3

import (
	. "github.com/trapped/gomaild/processors/pop3/client"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor"
	"log"
	"net"
	"os"
	"strconv"
	//POP3 commands
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/list"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/pass"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/stat"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/user"
)

type POP3 struct {
	//Port to listen at
	Port int
	//Whether to keep accepting clients (doesn't prevent active clients from continuing their current sessions)
	Keep bool
}

func (p *POP3) Listen() {
	log.Println("POP3: Starting POP3 server")
	if p.Keep == false {
		p.Keep = true
	}
	if p.Port == 0 {
		p.Port = 110
	}
	//Initialize POP3 commands
	cmdprocessor.Commands["user"] = user.Process
	cmdprocessor.Commands["pass"] = pass.Process
	cmdprocessor.Commands["stat"] = stat.Process
	cmdprocessor.Commands["list"] = list.Process

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
