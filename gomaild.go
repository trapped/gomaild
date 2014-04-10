package main

import (
	"fmt"
	"github.com/trapped/gomaild/config"
	"github.com/trapped/gomaild/processors/pop3"
	"github.com/trapped/gomaild/processors/smtp"
	"log"
	"runtime"
)

func main() {
	//Set max number of processors used to the number of CPUs
	runtime.GOMAXPROCS(runtime.NumCPU())

	config.Read()
	log.Println("gomaild: Starting gomaild")
	//Start POP3 server
	_pop3 := pop3.POP3{Port: 110, Keep: true}
	go _pop3.Listen()
	//Start SMTP server
	_smtp := smtp.SMTP{Port: 25, Keep: true}
	go _smtp.Listen()
	for {
		cmd := ""
		fmt.Scanln(&cmd)
		if cmd == "q" {
			break
		} else if cmd == "rc" {
			log.Println("gomaild: Reloading configuration")
			config.Read()
			log.Println("gomaild: Reloaded configuration")
		}
	}
}
