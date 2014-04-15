//Package textual provides parsing for telnet-like protocols.
package textual

import (
	"log"
	"math"
	"strings"
)

//A parser can be configured to remove a prefix, a suffix, heading and trailing whitespaces, split the text in arguments using separators, detect escaped characters by recognizing "brackets" (enclosures), and remove said brackets from the argument.
type Parser struct {
	Prefix             string //Prefix to remove
	Suffix             string //Suffix to remove
	OpenBrackets       bool   //Whether or not to open enclosures by removing the bracket characters
	Brackets           []byte //The characters to treat as brackets
	ArgumentSeparators []byte //The characters to treat as separators between arguments and therefore split the text by
	Trim               bool   //Whether or not to remove heading and trailing whitespaces.
}

//The object storing the result of the parse operation.
type Statement struct {
	Raw       string   //The raw input string
	Name      string   //The first argument (#0), which is usually the name of the command
	Arguments []string //The full list of arguments
}

//Parses the input string and outputs a Statement.
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

//Checks if a byte array contains a byte.
func inarray(a []byte, b byte) bool {
	for _, v := range a {
		if b == v {
			return true
		}
	}
	return false
}
