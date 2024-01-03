package copier

import (
	"reflect"
)

func copySQLScanner(dst reflect.Value, src reflect.Value) (err error) {
	if src.Type().AssignableTo(dst.Type()) {
		dst.Set(src)
		return
	}
	var valueField reflect.Value
	validField := dst.FieldByName("Valid")
	fields := fieldsOfStruct(dst.Type())
	for _, field := range fields {
		if field.Name == "Valid" {
			continue
		}
		valueField = dst.FieldByName(field.Name)
	}
	if err = copyValue(valueField, src); err != nil {
		return
	}
	if !valueField.IsZero() {
		validField.SetBool(true)
	}
	return
}
