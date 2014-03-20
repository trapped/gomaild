package session

type S_State int

var (
	AUTHORIZATION S_State = 0
	TRANSACTION   S_State = 1
	UPDATE        S_State = 2
)

type Session struct {
	RemoteEP      string
	State         S_State
	Authenticated bool
	Username      string
	Password      string
	Quitted       bool
	Retrieved     []interface{}
	Deleted       []interface{}
}
