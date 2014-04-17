//Package session provides structs to store session data.
package session

var (
	AUTHORIZATION int = 0 //Just connected
	TRANSACTION   int = 1 //Transaction in process: the client can read emails
	UPDATE        int = 2 //Client disconnecting
)

var (
	AUTHNONE      int = 0
	AUTHWUSER     int = 1 //AUTH Wait USERname
	AUTHWPASS     int = 2 //AUTH Wait PASSword
	AUTHENTICATED int = 3
)

type Session struct {
	RemoteEP      string        //Stores the string representation of the remote endpoint
	State         int           //Stores the session state
	InTLS         bool          //Whether or not the connection is being ecnrypted with TLS
	Authenticated bool          //Whether or not the client has authenticated
	AuthState     int           //Stores the authentication process state
	AuthMode      string        //Stores the authentication process mode
	Username      string        //Stores the username of the eventually authenticated client
	Password      string        //Stores the password of the eventually authenticated client
	Quitted       bool          //Whether or not the client has QUIT'd the connection
	Retrieved     []interface{} //Stores the messages to move in the "read" folder
	Deleted       []interface{} //Stores the messages to move to the "deleted" folder or delete
	Shared        string        //The session shared
}
