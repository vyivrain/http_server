package main

import "reflect"

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
