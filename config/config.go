//Package config provides configuration and templating.
package config

import (
	jconf "github.com/trapped/gomaild/parsers/config"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var Configuration Settings //Object to use when accessing configuration

//Root level configuration struct, storing global settings meant for the main package (the executable itself).
type Settings struct {
	Debug    bool              `json:"debug"`   //Whether or not to log debug information
	Accounts map[string]string `json:"Accounts` //User accounts
	Aliases  map[string]string `json:"Aliases"` //Email aliases
	POP3     POP3sett          `json:"POP3"`    //POP3 settings object
	SMTP     SMTPsett          `json:"SMTP"`    //SMTP settings object
	TLS      TLSsett           `json:"TLS"`     //TLS settings object
}

//Object storing POP3 settings.
type POP3sett struct {
	StartGreeting          string `json:"start_greeting"`           //Greeting message to use on connection start
	EndGreeting            string `json:"end_greeting"`             //Greeting message to use on connection end
	EnableUSER             bool   `json:"enable_user_cmd"`          //Whether or not to enable the USER/PASS login method
	SecureUSER             bool   `json:"secure_user_cmd"`          //Whether or not to accept invalid users to battle  trial-and-error bruteforce attacks (disables UserInvalidMessage)
	FakeDELE               bool   `json:"fake_dele_cmd"`            //Whether or not to move emails to the "deleted" folder instead of deleting them (makes them invisible to POP3 clients)
	Timeout                uint   `json:"timeout"`                  //Time in seconds between commands before timeout
	TimeoutMessage         string `json:"timeout_message"`          //Message to send on client timeout
	UserInvalidMessage     string `json:"user_invalid_message"`     //Message to send after an incorrect USER command (disabled by SecureUSER)
	UserOkMessage          string `json:"user_ok_message"`          //Message to send after a correct or "maybe" (see SecureUSER) USER command
	PasswordInvalidMessage string `json:"password_invalid_message"` //Message to send after an incorrect PASS/APOP command
	PasswordOkMessage      string `json:"password_ok_message"`      //Message to send after a correct PASS/APOP command
}

//Object storing SMTP settings.
type SMTPsett struct {
	StartGreeting           string `json:"start_greeting"`            //Greeting message to use on connection start
	EndGreeting             string `json:"end_greeting"`              //Greeting message to use on connection end
	QueuedMessage           string `json:"queued_message"`            //Message to send after an email has been successfully queued
	HelloMessage            string `json:"helo_message,ehlo_message"` //Greeting message to send along with HELO/EHLO replies
	SenderOkMessage         string `json:"sender_ok_message"`         //Message to send after a successful MAIL command
	SenderInvalidMessage    string `json:"sender_invalid_message"`    //Message to send after an unsuccessful MAIL command
	RecipientOkMessage      string `json:"recipient_ok_message"`      //Message to send after a successful RCPT command
	RecipientInvalidMessage string `json:"recipient_invalid_message"` //Message to send after an unsuccessful RCPT command
	DATAStartMessage        string `json:"data_start_message"`        //Message to send after the client requests to send email data
	Timeout                 uint   `json:"timeout"`                   //Time in seconds between commands before timeout
	TimeoutMessage          string `json:"timeout_message"`           //Message to send on client timeout
	EnableSTARTTLS          bool   `json:"enable_starttls"`           //Whether or not to enable the STARTTLS command
}

//Object storing TLS settings.
type TLSsett struct {
	CertificateFile    string `json:"certificate_file"`     //Path to a SSL certificate
	CertificateKeyFile string `json:"certificate_key_file"` //Path to a SSL certificate key
}

//Finds all *.conf files in the executable's directory and passes them to ParseConfig().
func Read() {
	currentfolder := path.Dir(os.Args[0])
	log.Println("Configuration: Searching files with pattern", currentfolder+"/*.conf")
	confs, err := filepath.Glob(currentfolder + "/*.conf")
	if err != nil {
		log.Println("Configuration: Error finding files:", err)
		return
	}
	log.Println("Configuration: Found", len(confs), "files")
	for _, v := range confs {
		ParseConfig(v)
	}
}

//Reads a file and parses its content into Configuration using the parsers/config package.
func ParseConfig(s string) {
	basename := strings.TrimSuffix(path.Base(s), ".conf")
	log.Println("Configuration: Parsing file", basename)
	file, err := ioutil.ReadFile(s)
	if err != nil {
		log.Println("Configuration: Error reading config file", basename+":", err)
		return
	}
	err = jconf.Parse(string(file), &Configuration)
	if err != nil {
		log.Println("Configuration: Error parsing config file", basename+":", err)
	}
}
