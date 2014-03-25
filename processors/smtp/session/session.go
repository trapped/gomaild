package session

type S_State int

var (
	NONE          S_State = 0
	IDENTIFICATED S_State = 1
	AUTHENTICATED S_State = 2
	RECAPITATION  S_State = 3
	COMPOSITION   S_State = 4
)

type Session struct {
	RemoteEP string
	State    S_State
	InSSL    bool
	Identity string
	Secret   string
	Quitted  bool
	Received []interface{}
}
