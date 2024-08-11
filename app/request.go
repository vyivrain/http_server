package main

import (
	"errors"
	"fmt"
	"slices"
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
	if hostIndex := r.matchRespectiveHeader(splittedMessage, "Host:"); hostIndex >= 0 {
		r.headers["host"] = strings.Replace(splittedMessage[hostIndex], "Host: ", "", -1)
	}
	if userAgentIndex := r.matchRespectiveHeader(splittedMessage, "User-Agent:"); userAgentIndex >= 0 {
		r.headers["userAgent"] = strings.Replace(splittedMessage[userAgentIndex], "User-Agent: ", "", -1)
	}
}

func (r *Request) matchEndpoint(appConfig *AppConfig) (*Endpoint, error) {
	for _, endpoint := range appConfig.endpoints {
		if endpoint.MatchesPath(r.headers["path"], r.headers["requestType"]) {
			return &endpoint, nil
		}
	}

	return nil, errors.New("")
}

func (r *Request) matchRespectiveHeader(headerSlice []string, header string) int {
	index := slices.IndexFunc(headerSlice, func(cmpHeader string) bool { return strings.Contains(cmpHeader, header) })
	return index
}
