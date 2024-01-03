package copier

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"reflect"
	"time"
)

var (
	timeType            = reflect.TypeOf(time.Time{})
	sqlScannerType      = reflect.TypeOf((*sql.Scanner)(nil)).Elem()
	sqlValuerType       = reflect.TypeOf((*driver.Valuer)(nil)).Elem()
	textMarshalerType   = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
	textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
)

func fieldsOfStruct(typ reflect.Type) (fields []reflect.StructField) {
	n := typ.NumField()
	for i := 0; i < n; i++ {
		field := typ.Field(i)
		if field.Anonymous {
			if field.Type.Kind() == reflect.Ptr && !field.IsExported() {
				continue
			}
			fields = append(fields, fieldsOfStruct(field.Type)...)
			continue
		}
		if !field.IsExported() {
			continue
		}
		fields = append(fields, field)
	}
	return
}
