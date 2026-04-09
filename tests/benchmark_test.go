package tests

import (
	"bufio"
	"net"
	"testing"
)

func BenchmarkSingleClientWrite(b *testing.B) {
	conn, err := net.Dial("tcp", "localhost:8090")
	if err != nil {
		b.Fatalf("Failed to connect to server. Is it running? Error: %v", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	cmd := []byte("*3\r\n$3\r\nSET\r\n$8\r\nbenchkey\r\n$3\r\n100\r\n")

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < 1; i++ {
		if _, err := conn.Write(cmd); err != nil {
			b.Fatalf("Failed to write to connection: %v", err)
		}

		// We assume your server replies with something ending in a newline (like "+OK\r\n")
		if _, err := reader.ReadBytes('\n'); err != nil {
			b.Fatalf("Failed to read response from server: %v", err)
		}
	}

	defer conn.Close()
	defer conn.Write([]byte("+QUIT\r\n"))

}
