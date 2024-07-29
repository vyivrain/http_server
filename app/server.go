package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

const CONN_TIMEOUT time.Duration = 10 * time.Second
const READ_TIMEOUT time.Duration = 100 * time.Millisecond

func readConnectionMessage(conn net.Conn) (string, error) {
	buffer := make([]byte, 2048)
	message := ""
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			// Timeout
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				break
			}
			// EOF
			if err == io.EOF {
				break
			}
		}
		data := buffer[:n]
		message += string(data)
	}

	return message, nil
}

func handleConnection(conn net.Conn, appConfig *AppConfig) {
	defer conn.Close()

	message, _ := readConnectionMessage(conn)
	request := Request{message: message}

	response, _ := request.Handle(appConfig)
	conn.Write([]byte(response))
}

func initAppConfigs() *AppConfig {
	router := Router{}

	return &AppConfig{
		endpoints: []Endpoint{
			{path: "/", requestType: "GET", handler: router.home, contentType: "text/plain"},
			{path: "/echo/{str}", requestType: "GET", handler: router.echo, contentType: "text/plain"},
		},
	}
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
	appConfig := initAppConfigs()

	for {
		conn, err := l.Accept()
		conn.SetDeadline(time.Now().Add(CONN_TIMEOUT))
		conn.SetReadDeadline(time.Now().Add(READ_TIMEOUT))

		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				fmt.Println("Timeout: ", err.Error())
			} else {
				fmt.Println("Error accepting connection: ", err.Error())
			}

			os.Exit(1)
		}

		go handleConnection(conn, appConfig)
	}
}
