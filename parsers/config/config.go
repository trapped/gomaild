package config

import (
	"encoding/json"
	"strings"
)

func Parse(text string, object interface{}) error {
	//Get original lines
	lines := strings.Split(text, "\n")
	//Remove junk spaces from lines
	for i, v := range lines {
		lines[i] = strings.TrimSpace(v)
	}
	//Remove completely commented or empty lines
	for i := 0; i < len(lines); i++ {
        v := lines[i]
		if strings.TrimSpace(v) == "" || v[0] == '#' {
			lines = append(lines[:i], lines[i+1:]...)
		}
	}
	//Trim away comments from lines
	for i, v := range lines {
		incomment := false
		cleaned := ""
		for c := 0; c < len(v); c++ {
			if v[c] == '#' {
				if incomment {
					incomment = false
				} else {
					incomment = true
				}
			}
			if !incomment {
				cleaned += string(v[c])
			}
		}
		lines[i] = cleaned
	}
	//Join lines again
	clean := strings.Join(lines, "\r\n")
	//Parse JSON text into object
	return json.Unmarshal([]byte(clean), &object)
}
