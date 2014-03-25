package error

import "strconv"

type Error struct {
	Code    int
	Message string
}

func (e *Error) New(c int, m string) Error {
	return Error{Code: c, Message: m}
}

func (e *Error) String() string {
	return strconv.Itoa(e.Code) + " " + e.Message
}
