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

func readFile(filePath string) (string, error) {
	if _, err := os.Stat(filePath); err != nil {
		return "", err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

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

func createFile(filePath string, data string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err2 := f.Write([]byte(data)); err2 != nil {
		return err2
	}

	return nil
}

func writeToFile(filePath string, data string) error {
	if _, err := os.Stat(filePath); err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err2 := f.Write([]byte(data)); err2 != nil {
		return err2
	}

	return nil
}

func deleteEmptyElements(slice interface{}) interface{} {
	sliceValue := reflect.ValueOf(slice)
	sliceLen := sliceValue.Len()
	newSlice := reflect.MakeSlice(sliceValue.Type(), 0, sliceValue.Len())
	if sliceValue.Kind() != reflect.Slice {
		panic("Input must be a slice")
	}

	for i := 0; i < sliceLen; i++ {
		value := sliceValue.Index(i)
		interfacedValue := value.Interface()
		switch convertedValue := interfacedValue.(type) {
		case string:
			if convertedValue != "" {
				newSlice = reflect.Append(newSlice, value)
			}
		default:
			if convertedValue != nil {
				newSlice = reflect.Append(newSlice, value)
			}
		}
	}

	return newSlice.Interface()
}
