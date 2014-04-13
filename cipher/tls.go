package cipher

import (
	"errors"
	"github.com/trapped/gomaild/config"
	. "github.com/trapped/gomaild/parsers/textual"
	"strings"
)

var TLS_CERT_FILE string
var TLS_CERT_KEY_FILE string
var TLS_CERT_SERVNAME string

func readconf() error {
}
