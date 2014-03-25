package smtp

import (
	. "github.com/trapped/gomaild/processors/smtp/client"
	"github.com/trapped/gomaild/processors/smtp/cmdprocessor"
	"log"
	"net"
	"os"
	"strconv"
	//SMTP commands
	//"github.com/trapped/gomaild/processors/smtp/cmdprocessor/data"
	"github.com/trapped/gomaild/processors/smtp/cmdprocessor/ehlo"
	"github.com/trapped/gomaild/processors/smtp/cmdprocessor/helo"
	"github.com/trapped/gomaild/processors/smtp/cmdprocessor/mail"
	"github.com/trapped/gomaild/processors/smtp/cmdprocessor/noop"
	"github.com/trapped/gomaild/processors/smtp/cmdprocessor/quit"
	"github.com/trapped/gomaild/processors/smtp/cmdprocessor/rcpt"
	//"github.com/trapped/gomaild/processors/smtp/cmdprocessor/rset"
	//"github.com/trapped/gomaild/processors/smtp/cmdprocessor/user"
	//Additional SMTP commands
	//"github.com/trapped/gomaild/processors/smtp/cmdprocessor/apop"
	//"github.com/trapped/gomaild/processors/smtp/cmdprocessor/capa"
	//"github.com/trapped/gomaild/processors/smtp/cmdprocessor/top"
	//"github.com/trapped/gomaild/processors/smtp/cmdprocessor/uidl"
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
	//Initialize SMTP commands
	cmdprocessor.Commands["helo"] = helo.Process
	cmdprocessor.Commands["ehlo"] = ehlo.Process
	cmdprocessor.Commands["quit"] = quit.Process
	cmdprocessor.Commands["mail"] = mail.Process
	cmdprocessor.Commands["rcpt"] = rcpt.Process
	//cmdprocessor.Commands["data"] = data.Process
	cmdprocessor.Commands["noop"] = noop.Process
	//cmdprocessor.Commands["quit"] = quit.Process
	//cmdprocessor.Commands["rset"] = rset.Process
	////Additional (non-compulsory in RFC1725) commands
	//cmdprocessor.Commands["uidl"] = uidl.Process
	//cmdprocessor.Commands["top"] = top.Process
	//cmdprocessor.Commands["apop"] = apop.Process
	//cmdprocessor.Commands["capa"] = apop.Process

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
