package main

import (
	"regexp"
	"strings"
)

type AppConfig struct {
	endpoints []Endpoint
}

type Endpoint struct {
	requestType string
	path        string
	handler     func(map[string]string) *Response
	params      map[string]string
	contentType string
}

func (e *Endpoint) MatchesPath(path string, requestType string) bool {
	if requestType != e.requestType {
		return false
	}

	if path == e.path {
		return true
	}

	requestedPathSplit := strings.Split(path, "/")
	comparePathSplit := strings.Split(e.path, "/")
	if len(requestedPathSplit) != len(comparePathSplit) {
		return false
	}

	zip := createZip(requestedPathSplit, comparePathSplit)
	for {
		requestPathValue, comparePathValue := zip()
		if requestPathValue == nil || comparePathValue == nil {
			break
		}

		if !e.matchParameterizedPath(requestPathValue.(string), comparePathValue.(string)) {
			return false
		}
	}

	return true
}

func (e *Endpoint) matchParameterizedPath(requestedSplittedPathValue string, compareSplittedPathValue string) bool {
	if strings.Contains(compareSplittedPathValue, "{str}") {
		re := regexp.MustCompile(`.+`)
		return re.MatchString(requestedSplittedPathValue)
	} else {
		return requestedSplittedPathValue == compareSplittedPathValue
	}
}

func (e *Endpoint) FillEndpointData(requestPath string) {
	e.fillParamsData(requestPath)
}

func (e *Endpoint) fillParamsData(requestPath string) {
	params := make(map[string]string)
	requestedPathSplit := strings.Split(requestPath, "/")
	comparePathSplit := strings.Split(e.path, "/")
	zip := createZip(requestedPathSplit, comparePathSplit)
	re := regexp.MustCompile(`[{}]`)

	for {
		requestPathValue, comparePathValue := zip()
		if requestPathValue == nil || comparePathValue == nil {
			break
		}

		strCompareValue := comparePathValue.(string)
		strRequestPathValue := requestPathValue.(string)

		if strings.Contains(strCompareValue, "{") {
			paramsKey := re.ReplaceAll([]byte(strCompareValue), []byte(""))
			params[string(paramsKey)] = strRequestPathValue
		}
	}

	e.params = params
}

func (e *Endpoint) GenerateResponse(httpVersion string, contentType string) *Response {
	response := e.handler(e.params)
	response.httpVersion = httpVersion
	response.contentType = contentType
	return response
}
