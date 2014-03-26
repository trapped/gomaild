package textual

import (
	"log"
	"math"
	"strings"
)

type Parser struct {
	Prefix             string
	Suffix             string
	OpenBrackets       bool
	Brackets           []byte
	ArgumentSeparators []byte
	Trim               bool
}

type Statement struct {
	Raw       string
	Name      string
	Arguments []string
}

func (p *Parser) Parse(s string) Statement {
	cmd := Statement{Raw: s}

	if p.Trim {
		s = strings.TrimSpace(s)
	}

	s = strings.TrimPrefix(s, p.Prefix)
	s = strings.TrimSuffix(s, p.Suffix)

	temparg := ""
	bracksopen := false

	for i := 0; i < len(s); i++ {
		//if (NOT IN ArgumentSeparators OR bracksopen) AND (NOT IN Brackets OR NOT OpenBrackets)
		if (!inarray(p.ArgumentSeparators, s[i]) || bracksopen) && (!inarray(p.Brackets, s[i]) || !p.OpenBrackets) {
			temparg += string(s[i])
		}
		if inarray(p.ArgumentSeparators, s[i]) && !bracksopen { //if (IN ArgumentSeparators) AND (NOT bracksopen)
			if temparg != "" {
				cmd.Arguments = append(cmd.Arguments, temparg)
				temparg = ""
			}
			continue
		} else if inarray(p.Brackets[:int(math.Ceil(float64(len(p.Brackets))/2))], s[i]) && !bracksopen { //if (IN FIRST HALF OF Brackets) AND (NOT bracksopen)
			bracksopen = true
			if temparg != "" {
				cmd.Arguments = append(cmd.Arguments, temparg)
				temparg = ""
			}
			continue
		} else if inarray(p.Brackets[int(math.Floor(float64(len(p.Brackets))/2)):], s[i]) && bracksopen { //if (IN SECOND HALF OF Brackets) AND (bracksopen)
			bracksopen = false
			if temparg != "" {
				cmd.Arguments = append(cmd.Arguments, temparg)
				temparg = ""
			}
			continue
		}
	}
	if temparg != "" {
		cmd.Arguments = append(cmd.Arguments, temparg)
	}

	if len(cmd.Arguments) > 0 {
		cmd.Name = cmd.Arguments[0]
	}

	d := func(f Statement) string {
		result := ""
		for i := 0; i < len(f.Arguments); i++ {
			result += "<" + f.Arguments[i] + ">"
		}
		return result
	}(cmd)
	log.Println(d)

	return cmd
}

func inarray(a []byte, b byte) bool {
	for _, v := range a {
		if b == v {
			return true
		}
	}
	return false
}
