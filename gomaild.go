package main

import (
	"fmt"
	"github.com/trapped/gomaild/config"
	"github.com/trapped/gomaild/processors/pop3"
	"log"
)

func main() {
	config.Read()
	fmt.Println("Current settings:", config.Settings)
	log.Println("Starting gomaild")
	//Start POP3 server
	_pop3 := pop3.POP3{Port: 110, Keep: true}
	go _pop3.Listen()
	for {
		cmd := ""
		fmt.Scanln(&cmd)
		if cmd == "q" {
			break
		} else if cmd == "rc" {
			config.Read()
			fmt.Println("Current settings:", config.Settings)
		}
	}
}
