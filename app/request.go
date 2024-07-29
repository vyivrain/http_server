package main

import (
	"errors"
	"fmt"
	"strings"
)

const NOT_FOUND_MESSAGE string = "Not Found"

type Request struct {
	message string
	headers map[string]string
}

func (r *Request) String() string {
	return fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n",
		r.headers["requestType"],
		r.headers["path"],
		r.headers["httpVersion"],
		r.headers["host"],
		r.headers["userAgent"],
	)
}

func (r *Request) Handle(appConfig *AppConfig) (string, error) {
	r.fillRequestData()
	endpoint, err := r.matchEndpoint(appConfig)
	if err != nil {
		response := &Response{
			message:    err.Error(),
			statusCode: 404,
			headers:    r.headers,
		}
		return fmt.Sprintf("%v", response), nil
	}

	return fmt.Sprintf("%v", endpoint.GenerateResponse(r.headers)), nil
}

func (r *Request) fillRequestData() {
	r.headers = make(map[string]string)
	splittedMessage := strings.Split(r.message, "\r\n")
	mainRequestInfo := strings.Split(splittedMessage[0], " ")

	r.headers["requestType"] = mainRequestInfo[0]
	r.headers["path"] = mainRequestInfo[1]
	r.headers["httpVersion"] = mainRequestInfo[2]
	r.headers["host"] = strings.Replace(splittedMessage[1], "Host: ", "", -1)
	r.headers["userAgent"] = strings.Replace(splittedMessage[2], "User-Agent: ", "", -1)
	r.headers["contentType"] = "text/plain"
}

func (r *Request) matchEndpoint(appConfig *AppConfig) (*Endpoint, error) {
	for _, endpoint := range appConfig.endpoints {
		if endpoint.MatchesPath(r.headers["path"], r.headers["requestType"]) {
			return &endpoint, nil
		}
	}

	return nil, errors.New("")
}
