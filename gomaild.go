package gomaild

import (
	"github.com/trapped/gomaild/processors/pop3"
	"log"
)

func main() {
	log.Println("Starting gomaild")
	go pop3.Listen()
}
