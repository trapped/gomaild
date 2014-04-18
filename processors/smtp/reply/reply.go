//Provides a data structure for SMTP replies.
package reply

import "strconv"

type Reply struct {
	Code    int
	Message string
}

func (r *Reply) String() string {
	return strconv.Itoa(r.Code) + " " + r.Message
}
