package textual

import (
	"strings"
)

type Parser struct {
	Prefix            string
	Suffix            string
	OpenBrackets      bool
	Brackets          string
	ArgumentSeparator byte
	MaxCmdLength      int
	Trim              bool
}

type Command struct {
	Raw       string
	Name      string
	Arguments []string
}

func (p *Parser) Parse(s string) Command {
	cmd := Command{
		Raw: s,
	}
	if p.Trim {
		s = strings.TrimSpace(s)
	}
	if p.Prefix != "" && strings.HasPrefix(s, p.Prefix) {
		s = strings.TrimPrefix(s, p.Prefix)
	}
	if p.Suffix != "" && strings.HasSuffix(s, p.Suffix) {
		s = strings.TrimSuffix(s, p.Suffix)
	}
	if p.Brackets != "" {
		arg := ""
		ob := false
		for i := 0; i < len(s); i++ {
			z := s[i]
			if z == p.Brackets[0] {
				ob = true
				if !p.OpenBrackets {
					arg += string(z)
				}
			} else if z == p.Brackets[1] {
				ob = false
				if !p.OpenBrackets {
					arg += string(z)
				}
			} else {
				if !ob {
					if z != p.ArgumentSeparator {
						fmt.Println("adding")
						arg += string(z)
					} else {
						if arg != "" {
							fmt.Println("appending")
							cmd.Arguments = append(cmd.Arguments, arg)
							arg = ""
						}
					}
				} else {
					fmt.Println("adding")
					arg += string(z)
				}
			}
			if i == len(s)-1 {
				if arg != "" {
					fmt.Println("appending")
					cmd.Arguments = append(cmd.Arguments, arg)
				}
			}
		}
	} else {
		cmd.Arguments = strings.Split(s, string(p.ArgumentSeparator))
	}
	if p.MaxCmdLength != -1 {
		cmd.Name = s[0:p.MaxCmdLength]
	} else {
		cmd.Name = cmd.Arguments[0]
	}

	return cmd
}
