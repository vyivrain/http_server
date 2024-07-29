package main

type Router struct{}

func (route *Router) home(params map[string]string) *Response {
	return &Response{statusCode: 200, message: ""}
}

func (route *Router) echo(params map[string]string) *Response {
	return &Response{statusCode: 200, message: params["str"]}
}
