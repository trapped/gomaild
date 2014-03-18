package textual

import (
	"strings"
)

type Parser struct {
	Prefix             string
	Suffix             string
	OpenBrackets       bool
	Brackets           string
	ArgumentSeparators []byte
	MaxCmdLength       int
	Trim               bool
	IgnoreEmpty        bool
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
		inbrackets := false
		considerwhites := false
		for i := 0; i < len(s); i++ {
			z := s[i]
			if z == p.Brackets[0] && inbrackets == false {
				inbrackets = true
				considerwhites = true
				if !p.OpenBrackets {
					arg += string(z)
				}
			} else if z == p.Brackets[1] && inbrackets == true {
				inbrackets = false
				considerwhites = false
				if !p.OpenBrackets {
					arg += string(z)
				}
			} else {
				if !inbrackets {
					if !inarray(p.ArgumentSeparators, z) {
						arg += string(z)
					} else {
						if (strings.TrimSpace(arg) != "" || !p.IgnoreEmpty || considerwhites) && arg != "" {
							cmd.Arguments = append(cmd.Arguments, arg)
							arg = ""
						}
					}
				} else {
					arg += string(z)
				}
			}
			if i == len(s)-1 {
				if (strings.TrimSpace(arg) != "" || !p.IgnoreEmpty || considerwhites) && arg != "" {
					cmd.Arguments = append(cmd.Arguments, arg)
					considerwhites = false
				}
			}
		}
	} else {
		cmd.Arguments = strings.Split(s, string(p.ArgumentSeparators[0]))
		if p.IgnoreEmpty {
			for i, v := range cmd.Arguments {
				if strings.TrimSpace(v) == "" {
					cmd.Arguments = append(cmd.Arguments[:i], cmd.Arguments[i+1:]...)
				}
			}
		}
	}
	if p.MaxCmdLength > 0 {
		cmd.Name = s[0:p.MaxCmdLength]
	} else {
		cmd.Name = cmd.Arguments[0]
	}

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
