package codeGenerator

import (
	"fmt"
	"reflect"
)

type Param struct {
	Math
	Keras
}

// For check there is any empty fields in Param
func checkNil(object interface{}) error {
	// Parameter object MUST BE POINTER STRUCT

	errorString := ""

	e := reflect.ValueOf(object).Elem()
	n := e.NumField()
	for i := 0; i < n; i++ {
		value := e.Field(i)
		tType := e.Type()

		// append error which field is nil
		if reflect.ValueOf(value.Interface()).IsNil() {
			errorString += fmt.Sprintf("field %s is nil\n", tType.Field(i).Name)
		}
	}

	if errorString == "" {
		return nil
	}

	return fmt.Errorf(errorString)
}
