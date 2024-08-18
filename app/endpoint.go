package main

import (
	"maps"
	"regexp"
	"strings"
)

type Endpoint struct {
	headers map[string]string
	handler func(map[string]string, map[string]string) *Response
}

func (e *Endpoint) MatchesPath(path string, requestType string) bool {
	if requestType != e.headers["requestType"] {
		return false
	}

	if path == e.headers["path"] {
		return true
	}

	requestedPathSplit := strings.Split(path, "/")
	comparePathSplit := strings.Split(e.headers["path"], "/")
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
	if strings.Contains(compareSplittedPathValue, "{") && strings.Contains(compareSplittedPathValue, "}") {
		re := regexp.MustCompile(`.+`)
		return re.MatchString(requestedSplittedPathValue)
	} else {
		return requestedSplittedPathValue == compareSplittedPathValue
	}
}

func (e *Endpoint) getParams(requestPath string) map[string]string {
	params := make(map[string]string)
	requestedPathSplit := strings.Split(requestPath, "/")
	comparePathSplit := strings.Split(e.headers["path"], "/")
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

	return params
}

func (e *Endpoint) GenerateResponse(headers map[string]string, params map[string]string) *Response {
	mergedHeaders := make(map[string]string)
	maps.Copy(mergedHeaders, e.headers)
	maps.Copy(mergedHeaders, headers)
	mergedParams := make(map[string]string)
	maps.Copy(mergedParams, e.getParams(mergedHeaders["path"]))
	maps.Copy(mergedParams, params)
	response := e.handler(mergedHeaders, mergedParams)
	maps.Copy(mergedHeaders, response.headers)
	response.headers = mergedHeaders

	return response
}
