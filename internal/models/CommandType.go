package models

type CommandType int

const (
	PING = iota
	SET
	GET
	UNKNOWN
)

var commandTypeToString = map[CommandType]string{
	PING:    "PING",
	SET:     "SET",
	GET:     "GET",
	UNKNOWN: "UNKNOWN",
}

var CommandStringToType = map[string]CommandType{
	"PING":    PING,
	"SET":     SET,
	"GET":     GET,
	"UNKNOWN": UNKNOWN,
}

func (c CommandType) String() string {
	if str, ok := commandTypeToString[c]; ok {
		return str
	}
	return "UNKNOWN COMMAND"
}
