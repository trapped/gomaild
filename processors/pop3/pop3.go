package pop3

import (
	. "github.com/trapped/gomaild/processors/pop3/client"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var (
	PORT int  = 110
	KEEP bool = true
)

func Listen() {
	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(PORT))
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	for KEEP {
		client, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		cliobj := MakeClient(&listener, client)
		go cliobj.proc()
	}
}
