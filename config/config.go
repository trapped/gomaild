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

var Configuration Settings

type Settings struct {
	Debug    bool              `json:"debug"`
	Accounts map[string]string `json:"Accounts`
	Aliases  map[string]string `json:"Aliases"`
	POP3     POP3sett          `json:"POP3"`
	SMTP     SMTPsett          `json:"SMTP"`
	TLS      TLSsett           `json:"TLS"`
}

type POP3sett struct {
	StartGreeting string `json:"start_greeting"`
	EndGreeting   string `json:"end_greeting"`
	EnableUSER    bool   `json:"enable_user_cmd"`
	SecureUSER    bool   `json:"secure_user_cmd"`
	FakeDELE      bool   `json:"fake_dele_cmd"`
}

type SMTPsett struct {
	StartGreeting           string `json:"start_greeting"`
	EndGreeting             string `json:"end_greeting"`
	QueuedMessage           string `json:"queued_message"`
	HelloMessage            string `json:"helo_message,ehlo_message"`
	SenderOkMessage         string `json:"sender_ok_message"`
	SenderInvalidMessage    string `json:"sender_invalid_message"`
	RecipientOkMessage      string `json:"recipient_ok_message"`
	RecipientInvalidMessage string `json:"recipient_invalid_message"`
	DATAStartMessage        string `json:"data_start_message"`
	Timeout                 uint   `json:"timeout"`
	TimeoutMessage          string `json:"timeout_message"`
}

type TLSsett struct {
	CertificateFile    string `json:"certificate_file"`
	CertificateKeyFile string `json:"certificate_key_file"`
}

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
