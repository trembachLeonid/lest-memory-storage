package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const addr = "localhost:8090"

func main() {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect %s: %v\n", addr, err)
		os.Exit(1)
	}
	defer conn.Close()

	r := bufio.NewReader(conn)
	send := func(args ...string) string {
		return sendCmd(conn, r, args...)
	}

	// run("PING", send("PING", ""))

	run("SET bar bar", send("SET", "bar", "bar"))
	run("EXPIRE bar 15", send("EXPIRE", "bar", "21"))
	run("GET bar", send("GET", "bar"))
	run("SET foo bar", send("SET", "foo", "bar"))
	run("EXPIRE foo 10", send("EXPIRE", "foo", "10"))
	run("GET foo", send("GET", "foo"))
	// time.Sleep(11 * time.Second)
	// run("GET foo", send("GET", "foo"))
	// run("SET counter 10", send("SET", "counter", "10"))
	// run("INC counter", send("INC", "counter", ""))
	// run("INC counter 5", send("INC", "counter", "5"))
	// run("DEC counter 3", send("DEC", "counter", "3"))
	// run("GET counter", send("GET", "counter"))
	// run("DEL foo", send("DEL", "foo", ""))
	// run("GET foo (deleted)", send("GET", "foo"))
	// run("CONFIG", send("CONFIG", ""))
	run("QUIT", send("QUIT", ""))
}

// sendCmd encodes args as a RESP array and returns the server response line.
func sendCmd(conn net.Conn, r *bufio.Reader, args ...string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "*%d\r\n", len(args))
	for _, a := range args {
		fmt.Fprintf(&sb, "$%d\r\n%s\r\n", len(a), a)
	}
	if _, err := fmt.Fprint(conn, sb.String()); err != nil {
		return fmt.Sprintf("write error: %v", err)
	}
	line, err := r.ReadString('\n')
	if err != nil {
		return fmt.Sprintf("read error: %v", err)
	}
	return strings.TrimRight(line, "\r\n")
}

func run(label, response string) {
	fmt.Printf("%-28s => %s\n", label, response)
}
