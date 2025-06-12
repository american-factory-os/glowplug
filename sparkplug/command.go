package sparkplug

// Command is the type of message being sent
type Command string

// Command is a NODE or DEVICE command
const (
	// device birth event
	DBIRTH Command = "DBIRTH"
	// device data event
	DDATA Command = "DDATA"
	// device command event
	DCMD Command = "DCMD"
	// device death event
	DDEATH Command = "DDEATH"
	// node command event
	NCMD Command = "NCMD"
	// node birth event
	NBIRTH Command = "NBIRTH"
	// node data event
	NDATA Command = "NDATA"
	// node death event
	NDEATH Command = "NDEATH"
	// critical application state message
	STATE Command = "STATE"
)

const (
	// The Will Payload will be the UTF-8 STRING “OFFLINE”.
	OFFLINE = "OFFLINE"
	// The Birth Certificate Payload is the UTF-8 STRING “ONLINE”
	ONLINE = "ONLINE"
)

// String returns the string representation of the Command
func (x Command) String() string {
	return string(x)
}
