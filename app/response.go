package main

import "fmt"

type Response struct {
	message    string
	statusCode int
	headers    map[string]string
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
		"Content-Type: %s\r\nContent-Length: %d\r\n\r\n",
		resp.headers["contentType"],
		len(resp.message),
	)

	return htmlStatus + htmlHeaders + resp.message
}

func (resp *Response) getHeadMessage() string {
	switch resp.statusCode {
	case 200:
		return "OK"
	case 404:
		return "Not Found"
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
