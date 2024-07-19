package main

import (
	"errors"
	"fmt"
	"strings"
)

const NOT_FOUND_MESSAGE string = "Not Found"

type Request struct {
	message     string
	requestType string
	path        string
	httpVersion string
	host        string
	userAgent   string
}

type Response struct {
	message     string
	statusCode  int
	httpVersion string
}

func (r *Request) String() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n", r.requestType, r.path, r.httpVersion, r.host, r.userAgent)
}

func (resp *Response) String() string {
	return fmt.Sprintf("%s %d %s\r\n\r\n", resp.httpVersion, resp.statusCode, resp.message)
}

func (r *Request) fillRequestData(message string) {
	splittedMessage := strings.Split(message, "\r\n")
	mainRequestInfo := strings.Split(splittedMessage[0], " ")

	r.message = message
	r.requestType = mainRequestInfo[0]
	r.path = mainRequestInfo[1]
	r.httpVersion = mainRequestInfo[2]
	r.host = strings.Replace(splittedMessage[1], "Host: ", "", -1)
	r.userAgent = strings.Replace(splittedMessage[2], "User-Agent: ", "", -1)
}

func (r *Request) Handle(message string) (string, error) {
	r.fillRequestData(message)
	err := r.ValidatePathNotFound()
	if err != nil {
		return fmt.Sprintf("%v", r.GenerateResponse(err.Error(), 404)), nil
	}

	switch r.requestType {
	case "GET":
		return fmt.Sprintf("%v", r.HandleGetRequest()), nil
	default:
		return fmt.Sprintf("%v", r.GenerateResponse(NOT_FOUND_MESSAGE, 404)), nil
	}
}

func (r *Request) HandleGetRequest() *Response {
	return r.GenerateResponse("OK", 200)
}

func (r *Request) ValidatePathNotFound() error {
	if r.path != "/" {
		return errors.New(NOT_FOUND_MESSAGE)
	} else {
		return nil
	}
}

func (r *Request) GenerateResponse(message string, statusCode int) *Response {
	return &Response{statusCode: statusCode, message: message, httpVersion: r.httpVersion}
}
