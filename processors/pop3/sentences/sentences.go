//Package sentences defines a set of function to parse, elaborate and get POP3 sentences (such as greetings) from the configuration file.
package sentences

import (
	"github.com/trapped/gomaild/config"
	. "github.com/trapped/gomaild/parsers/textual"
	"strings"
)

//StartGreeting gets the user-defined POP3 client greeting.
func StartGreeting() string {
	greeting := ""
	if config.Settings["pop3"] != nil && len(config.Settings["pop3"]["start_greeting"]) >= 1 && strings.TrimSpace(config.Settings["pop3"]["start_greeting"][0].(Statement).Arguments[1]) != "" {
		greeting += config.Settings["pop3"]["start_greeting"][0].(Statement).Arguments[1]
	}
	return greeting
}

//EndGreeting gets the user-defined POP3 client greeting.
func EndGreeting() string {
	greeting := ""
	if config.Settings["pop3"] != nil && len(config.Settings["pop3"]["end_greeting"]) >= 1 && strings.TrimSpace(config.Settings["pop3"]["end_greeting"][0].(Statement).Arguments[1]) != "" {
		greeting += config.Settings["pop3"]["end_greeting"][0].(Statement).Arguments[1]
	}
	return greeting
}
