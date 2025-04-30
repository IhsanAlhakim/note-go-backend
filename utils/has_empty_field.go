package utils

import (
	"net/http"
	"reflect"
)

func HasEmptyField(w http.ResponseWriter, s interface{}) bool {

	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			JSONResponse(w, R{Message: "Missing Credentials"}, http.StatusBadRequest)
			return true
		}
	}
	return false
}
