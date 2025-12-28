package tc

/*
Textcat Websocket Protocol (TWP)
*/


type Recieve struct {
	/* general fields */
	Req string // request type
	Username string
	Token string // also used as password when loggin in or registering
	Value string // a value for storing something, changes depending on request type

	/* request specific fields */
}

type Send struct {
	Req string // request type
	Value string
	Status string
	/*
	ok: went well
	validate_fail: input validation error
	server_error: internal server error
	*/
}


// bridge between the cmd/ and the tc/
type Handler interface {
	HandleReq(msg []byte)
	LogMsg(level string, message string, args ...any)
	Store(table string, record any)
}