//Package session provides structs to store session data.
package session

var (
	NONE         int = 0 //Just connected
	IDENTIFIED   int = 1 //EHLO/HELO has been issued already
	RECAPITATION int = 2 //MAIL has been issued with success
	COMPOSITION  int = 3 //DATA has been issued with success and it's receiving data
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
	InTLS         bool          //Whether or not the connection is being encrypted with TLS
	Authenticated bool          //Whether or not the client has authenticated
	AuthState     int           //Stores the authentication process state
	AuthMode      string        //Stores the authentication process mode
	Username      string        //Stores the username of the eventually authenticated client
	Password      string        //Stores the password of the eventually authenticated client
	Identity      string        //Stores the identity provided by the client with EHLO/HELO
	Shared        string        //Shared secret for CRAM-MD5-like authentication methods
	Quitted       bool          //Whether or not the client has QUIT'd the connection
	Received      []interface{} //Stores the messages to recapitate
}
