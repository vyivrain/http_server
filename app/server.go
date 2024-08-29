package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"slices"
	"time"
)

const connTimeout time.Duration = 30 * time.Second
const connChunkSize int = 4092

// const READ_TIMEOUT time.Duration = 200 * time.Millisecond

func readConnectionMessage(conn net.Conn) (string, error) {
	buffer := bytes.NewBuffer(nil)
	for {
		chunk := make([]byte, connChunkSize)
		read, err := conn.Read(chunk)
		if err != nil {
			return "", err
		}
		buffer.Write(chunk[:read])

		if read == 0 || read < connChunkSize {
			break
		}
	}

	return buffer.String(), nil
}

func handleConnection(conn net.Conn, appConfig *AppConfig) {
	defer conn.Close()

	message, _ := readConnectionMessage(conn)
	request := Request{message: message}

	response, _ := request.Handle(appConfig)
	conn.Write([]byte(response))
}

func initAppConfigs() *AppConfig {
	appConfig := AppConfig{}
	router := Router{}

	appConfig.endpoints = []Endpoint{
		{handler: router.home, headers: map[string]string{"requestType": "GET", "contentType": "text/plain", "path": "/"}},
		{handler: router.userAgent, headers: map[string]string{"requestType": "GET", "contentType": "text/plain", "path": "/user-agent"}},
		{handler: router.echo, headers: map[string]string{"requestType": "GET", "contentType": "text/plain", "path": "/echo/{str}"}},
		{handler: router.files, headers: map[string]string{"requestType": "GET", "contentType": "application/octet-stream", "path": "/files/{filename}"}},
		{handler: router.postFile, headers: map[string]string{"requestType": "POST", "contentType": "application/octet-stream", "path": "/files/{filename}"}},
	}

	args := os.Args[1:]
	directoryParamIndex := slices.Index(args, "--directory")
	if directoryParamIndex >= 0 {
		directory := args[directoryParamIndex+1]
		router.fileDirectory = directory

	}

	appConfig.router = router

	return &appConfig
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
		conn.SetDeadline(time.Now().Add(connTimeout))

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
