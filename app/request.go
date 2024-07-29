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

func (r *Request) String() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n", r.requestType, r.path, r.httpVersion, r.host, r.userAgent)
}

func (r *Request) Handle(appConfig *AppConfig) (string, error) {
	r.fillRequestData()
	endpoint, err := r.matchEndpoint(appConfig)
	if err != nil {
		response := &Response{message: err.Error(), statusCode: 404, httpVersion: r.httpVersion, contentType: "text/plain"}
		return fmt.Sprintf("%v", response), nil
	}
	endpoint.FillEndpointData(r.path)

	return fmt.Sprintf("%v", endpoint.GenerateResponse(r.httpVersion, "text/plain")), nil
}

func (r *Request) fillRequestData() {
	splittedMessage := strings.Split(r.message, "\r\n")
	mainRequestInfo := strings.Split(splittedMessage[0], " ")

	r.requestType = mainRequestInfo[0]
	r.path = mainRequestInfo[1]
	r.httpVersion = mainRequestInfo[2]
	r.host = strings.Replace(splittedMessage[1], "Host: ", "", -1)
	r.userAgent = strings.Replace(splittedMessage[2], "User-Agent: ", "", -1)
}

func (r *Request) matchEndpoint(appConfig *AppConfig) (*Endpoint, error) {
	for _, endpoint := range appConfig.endpoints {
		if endpoint.MatchesPath(r.path, r.requestType) {
			return &endpoint, nil
		}
	}

	return nil, errors.New("")
}
