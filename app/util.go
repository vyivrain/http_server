package main

import (
	"bytes"
	"io"
	"os"
	"reflect"
)

func createZip(slice1 interface{}, slice2 interface{}) func() (interface{}, interface{}) {
	s1 := reflect.ValueOf(slice1)
	s2 := reflect.ValueOf(slice2)

	minLen := s1.Len()

	if minLen > s2.Len() {
		minLen = s2.Len()
	}

	i := 0
	return func() (interface{}, interface{}) {
		i++
		if i > minLen {
			return nil, nil
		} else {
			return s1.Index(i - 1).Interface(), s2.Index(i - 1).Interface()
		}
	}
}

func readFile(filepath string) (string, error) {
	if _, err := os.Stat(filepath); err != nil {
		return "", err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}

	buffer := bytes.NewBuffer(nil)

	for {
		chunk := make([]byte, 1024)
		n, err := file.Read(chunk)
		if err != nil && err != io.EOF {
			panic(err)
		}

		if n == 0 {
			break
		}

		buffer.Write(chunk[:n])
	}

	return buffer.String(), nil
}
