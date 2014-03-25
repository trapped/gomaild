//Package sentences defines a set of function to parse, elaborate and get SMTP sentences (such as greetings) from the configuration file.
package sentences

import (
	"github.com/trapped/gomaild/config"
	. "github.com/trapped/gomaild/parsers/textual"
	"strings"
)

//StartGreeting gets the user-defined SMTP client greeting.
func StartGreeting() string {
	greeting := ""
	if config.Settings["smtp"] != nil && len(config.Settings["smtp"]["start_greeting"]) >= 1 && strings.TrimSpace(config.Settings["smtp"]["start_greeting"][0].(Statement).Arguments[1]) != "" {
		greeting += config.Settings["smtp"]["start_greeting"][0].(Statement).Arguments[1]
	}
	return greeting
}

//EndGreeting gets the user-defined SMTP client greeting.
func EndGreeting() string {
	greeting := ""
	if config.Settings["smtp"] != nil && len(config.Settings["smtp"]["end_greeting"]) >= 1 && strings.TrimSpace(config.Settings["smtp"]["end_greeting"][0].(Statement).Arguments[1]) != "" {
		greeting += config.Settings["smtp"]["end_greeting"][0].(Statement).Arguments[1]
	}
	return greeting
}
