package main

import "fmt"

type Response struct {
	message     string
	statusCode  int
	headers     map[string]string
	compression compressionMethod
}

type ResponseError struct {
	statusCode int
	message    string
}

func (resp *Response) String() string {
	htmlStatus := fmt.Sprintf("%s %d %s\r\n", resp.headers["httpVersion"], resp.statusCode, resp.getHeadMessage())
	if resp.statusCode == 404 {
		return htmlStatus + "\r\n"
	}

	htmlHeaders := fmt.Sprintf(
		"Content-Type: %s\r\nContent-Length: %d",
		resp.headers["contentType"],
		len(resp.message),
	)

	additionalHeaders := ""
	if val, ok := resp.headers["encoding"]; ok {
		additionalHeaders += "\r\nContent-Encoding: " + val
	}

	responseMessage := ""
	if resp.compression != nil {
		if compressedData, err := resp.compression.compress([]byte(resp.message)); err == nil {
			responseMessage = string(compressedData)
		} else {
			panic("Can't compress response message")
		}
	} else {
		responseMessage = resp.message
	}

	return htmlStatus + htmlHeaders + additionalHeaders + "\r\n\r\n" + responseMessage
}

func (resp *Response) getHeadMessage() string {
	switch resp.statusCode {
	case 200:
		return "OK"
	case 404:
		return "Not Found"
	case 201:
		return "Created"
	}

	return ""
}

func (respError *ResponseError) Error() string {
	if respError.statusCode == 404 || respError.message != "" {
		return ""
	} else {
		return "Unhandled Error"
	}
}
