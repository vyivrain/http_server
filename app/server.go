package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

const TIMEOUT_SECONDS time.Duration = 10 * time.Second

func readConnectionMessage(conn net.Conn) (string, error) {
	buf := make([]byte, 2048)
	read, _ := conn.Read(buf)
	fmt.Println(string(read))

	return string(read), nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	input, _ := readConnectionMessage(conn)
	fmt.Println(input)

	responseString := "HTTP/1.1 200 OK\r\n\r\n"
	conn.Write([]byte(responseString))
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		conn.SetDeadline(time.Now().Add(TIMEOUT_SECONDS))
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				fmt.Println("Timeout: ", err.Error())
			} else {
				fmt.Println("Error accepting connection: ", err.Error())
			}

			os.Exit(1)
		}

		go handleConnection(conn)
	}
}
