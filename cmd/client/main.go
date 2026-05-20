package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const addr = "localhost:8090"

var passed, failed int

func main() {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect %s: %v\n", addr, err)
		os.Exit(1)
	}
	defer conn.Close()

	r := bufio.NewReader(conn)
	send := func(args ...string) string { return sendCmd(conn, r, args...) }

	// SET / GET
	check("DEL todel", send("DEL", "todel"), "+OK")
	check("GET foo", send("GET", "foo"), "bar")

	// DEL
	check("SET todel x", send("SET", "todel", "x"), "+OK")
	check("GET todel (deleted)", send("GET", "todel"), "$-1")

	// GET missing key
	check("GET nonexistent", send("GET", "nonexistent"), "$-1")

	// INC / DEC
	check("SET counter 10", send("SET", "counter", "10"), "+OK")
	check("INC counter", send("INC", "counter"), "11")
	check("INC counter 5", send("INC", "counter", "5"), "16")
	check("DEC counter 3", send("DEC", "counter", "3"), "13")
	check("GET counter", send("GET", "counter"), "13")

	// EXPIRE
	check("SET temp value", send("SET", "temp", "value"), "+OK")
	check("EXPIRE temp 1s", send("EXPIRE", "temp", "1"), "+OK")
	check("GET temp (alive)", send("GET", "temp"), "value")
	fmt.Println("  waiting 2s for expiry...")
	time.Sleep(2 * time.Second)
	check("GET temp (expired)", send("GET", "temp"), "$-1")

	send("QUIT")
	fmt.Printf("\n%d passed, %d failed\n", passed, failed)
	if failed > 0 {
		os.Exit(1)
	}
}

func check(label, got, want string) {
	if got == want {
		fmt.Printf("  PASS  %-30s  got=%q\n", label, got)
		passed++
	} else {
		fmt.Printf("  FAIL  %-30s  got=%q  want=%q\n", label, got, want)
		failed++
	}
}

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
