package config

import (
	"github.com/trapped/gomaild/parsers/textual"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var Settings map[string]map[string][]interface{} = make(map[string]map[string][]interface{}, 0)

var parser = textual.Parser{
	OpenBrackets:      true,
	Brackets:          "``",
	Trim:              true,
	IgnoreEmpty:       true,
	ArgumentSeparator: ' ',
}

func Read() {
	currentfolder := path.Dir(os.Args[0])
	log.Println("Searching configuration files with pattern", currentfolder+"/*.conf")
	confs, err := filepath.Glob(currentfolder + "/*.conf")
	if err != nil {
		log.Println("Error finding configuration files:", err)
		return
	}
	log.Println("Found", len(confs), "configuration files")
	for _, v := range confs {
		ParseConfig(v)
	}
}

func ParseConfig(s string) {
	basename := strings.TrimSuffix(path.Base(s), ".conf")
	log.Println("Parsing configuration file", basename)
	file, err := ioutil.ReadFile(s)
	if err != nil {
		log.Println("Error reading config file", basename)
		return
	}
	filetext := string(file)
	{
		intocommentblock := false
		statement := ""
		for i := 0; i < len(string(file)); i++ {
			if filetext[i] == '#' {
				if intocommentblock {
					intocommentblock = false
				} else {
					intocommentblock = true
				}
			} else {
				if !intocommentblock {
					if filetext[i] != ';' {
						statement += string(filetext[i])
					} else {
						if strings.TrimSpace(statement) != "" {
							parsedstatement := parser.Parse(statement)
							if Settings[basename] == nil {
								Settings[basename] = make(map[string][]interface{})
							}
							Settings[basename][parsedstatement.Name] = append(Settings[basename][parsedstatement.Name], parsedstatement)
							statement = ""
						}
					}
				}
			}
		}
	}
}
