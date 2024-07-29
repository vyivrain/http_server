package main

type Router struct{}

func (route *Router) home(headers map[string]string, params map[string]string) *Response {
	return &Response{statusCode: 200, message: ""}
}

func (route *Router) echo(headers map[string]string, params map[string]string) *Response {
	return &Response{statusCode: 200, message: params["str"]}
}

func (route *Router) userAgent(headers map[string]string, params map[string]string) *Response {
	return &Response{statusCode: 200, message: headers["userAgent"]}
}
