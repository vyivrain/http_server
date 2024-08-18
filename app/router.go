package main

import "fmt"

type Router struct {
	fileDirectory string
}

func (route *Router) home(headers map[string]string, params map[string]string) *Response {
	return &Response{statusCode: 200, message: ""}
}

func (route *Router) echo(headers map[string]string, params map[string]string) *Response {
	return &Response{statusCode: 200, message: params["str"]}
}

func (route *Router) userAgent(headers map[string]string, params map[string]string) *Response {
	return &Response{statusCode: 200, message: headers["userAgent"]}
}

func (route *Router) files(headers map[string]string, params map[string]string) *Response {
	filePath := route.fileDirectory + params["filename"]

	fileData, err := readFile(filePath)
	if err != nil {
		fmt.Println(err)
		return &Response{statusCode: 404}
	}

	return &Response{statusCode: 200, message: fileData}
}

func (route *Router) postFile(headers map[string]string, params map[string]string) *Response {
	filePath := route.fileDirectory + params["filename"]
	_, err := readFile(filePath)

	if err != nil {
		if fileErr := createFile(filePath, params["body"]); fileErr != nil {
			fmt.Println(fileErr)
			return &Response{statusCode: 404}
		}
	} else {
		if fileErr := writeToFile(filePath, params["body"]); fileErr != nil {
			fmt.Println(fileErr)
			return &Response{statusCode: 404}
		}
	}

	fmt.Println(params)

	return &Response{statusCode: 201}
}
