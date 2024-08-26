package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type Request struct {
	message          string
	headers          map[string]string
	params           map[string]string
	body             string
	unhandledRequest bool
	compression      compressionMethod
}

func (r *Request) String() string {
	return fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n%s\n%t\n",
		r.headers["requestType"],
		r.headers["path"],
		r.headers["httpVersion"],
		r.headers["host"],
		r.headers["userAgent"],
		r.body,
		r.unhandledRequest,
	)
}

func (r *Request) Handle(appConfig *AppConfig) (string, error) {
	r.fillRequestData()
	endpoint, err := r.matchEndpoint(appConfig)
	if err != nil || r.unhandledRequest {
		response := &Response{
			message:    err.Error(),
			statusCode: 404,
			headers:    r.headers,
		}
		return fmt.Sprintf("%v", response), nil
	}

	if r.body != "" {
		if r.compression != nil && r.compression.dataContainsCompressionMethod([]byte(r.body)) {
			compressedBody, err2 := r.compression.decompress([]byte(r.body))
			if err2 != nil {
				panic("Can't compress request's body")
			}

			r.params["body"] = string(compressedBody)
		} else {
			r.params["body"] = r.body
		}
	}

	return fmt.Sprintf("%v", endpoint.GenerateResponse(r.headers, r.params, r.compression)), nil
}

func (r *Request) fillRequestData() {
	r.headers = make(map[string]string)
	r.params = make(map[string]string)
	splittedMessage := strings.Split(r.message, "\r\n")
	splittedMessage = deleteEmptyElements(splittedMessage).([]string)
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
	if contentTypeIndex := r.matchRespectiveHeader(splittedMessage, "Content-Type:"); contentTypeIndex >= 0 {
		r.headers["contentType"] = strings.Replace(splittedMessage[contentTypeIndex], "Content-Type: ", "", -1)
	}

	contentLengthIndex := r.matchRespectiveHeader(splittedMessage, "Content-Length:")
	if contentLengthIndex >= 0 {
		r.headers["contentLength"] = strings.Replace(splittedMessage[contentLengthIndex], "Content-Length: ", "", -1)
	}

	enncodingIndex := r.matchRespectiveHeader(splittedMessage, "Accept-Encoding:")
	if enncodingIndex >= 0 {
		encoding := strings.Replace(splittedMessage[enncodingIndex], "Accept-Encoding: ", "", -1)
		if r.validEncoding(encoding) {
			r.headers["encoding"] = encoding
		}
	}

	if slices.Contains([]string{"POST", "PATCH", "PUT"}, r.headers["requestType"]) {
		switch r.headers["contentType"] {
		case "application/octet-stream", "text/html", "text/plain":
			r.body = splittedMessage[len(splittedMessage)-1]
		default:
			r.unhandledRequest = true
		}
	}

	r.setEncoding()
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

func (r *Request) validEncoding(encoding string) bool {
	validEncodings := []string{"gzip"}
	return slices.Contains(validEncodings, encoding)
}

func (r *Request) setEncoding() {
	if r.headers["encoding"] == "gzip" {
		r.compression = &GzipCompression{}
	}
}
