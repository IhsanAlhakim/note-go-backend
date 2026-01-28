package validation

import (
	"errors"
	"reflect"
)

var ErrStructEmptyField = errors.New("empty field in struct")

func CheckStructFields(s any) error {

	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := range v.NumField() {
		field := v.Field(i)
		if reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			return ErrStructEmptyField
		}
	}
	return nil
}
