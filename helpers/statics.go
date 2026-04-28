package helpers

var QUIT = []byte("QUIT")
var PING = []byte("PING")
var PONG = []byte("+PONG")
var SET = []byte("SET")
var GET = []byte("GET")
var DEL = []byte("DEL")
var INC = []byte("INC")
var DEC = []byte("DEC")
var CONFIG = []byte("CONFIG")
var UNKNOWN_COMMAND = []byte("UNKNOWN COMMAND")
var OPERATION_ERROR = []byte("OPERATION ERROR")
var CRLF = []byte{'\r', '\n'}
var OK = []byte("+OK")
