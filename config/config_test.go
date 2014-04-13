package config

import (
	"testing"
)

func Test0(t *testing.T) {
	ParseConfig("/home/giorgio/go/src/github.com/trapped/gomaild/gomaild.conf")
	t.Log(Configuration)
}
