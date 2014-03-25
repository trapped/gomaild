package pop3

import (
	. "github.com/trapped/gomaild/processors/pop3/client"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor"
	"log"
	"net"
	"os"
	"strconv"
	//POP3 commands
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/dele"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/list"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/noop"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/pass"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/quit"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/retr"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/rset"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/stat"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/user"
	//Additional POP3 commands
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/apop"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/capa"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/top"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/uidl"
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
	cmdprocessor.Commands["retr"] = retr.Process
	cmdprocessor.Commands["dele"] = dele.Process
	cmdprocessor.Commands["noop"] = noop.Process
	cmdprocessor.Commands["quit"] = quit.Process
	cmdprocessor.Commands["rset"] = rset.Process
	//Additional (non-compulsory in RFC1725) commands
	cmdprocessor.Commands["uidl"] = uidl.Process
	cmdprocessor.Commands["top"] = top.Process
	cmdprocessor.Commands["apop"] = apop.Process
	cmdprocessor.Commands["capa"] = capa.Process

	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(p.Port))
	if err != nil {
		log.Println("POP3:", err)
		os.Exit(1)
	}
	for p.Keep {
		client, err := listener.Accept()
		if err != nil {
			log.Println("POP3:", err)
			continue
		}
		cliobj := MakeClient(&listener, client)
		log.Println("POP3: Accepting client", cliobj.RemoteEP())
		go cliobj.Process()
	}
}
