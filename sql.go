package copier

import (
	"database/sql"
	"reflect"
)

var (
	sqlNullStringType  = reflect.TypeOf(sql.NullString{})
	sqlNullInt16Type   = reflect.TypeOf(sql.NullInt16{})
	sqlNullInt32Type   = reflect.TypeOf(sql.NullInt32{})
	sqlNullInt64Type   = reflect.TypeOf(sql.NullInt64{})
	sqlNullFloat64Type = reflect.TypeOf(sql.NullFloat64{})
	sqlNullBoolType    = reflect.TypeOf(sql.NullBool{})
	sqlNullTimeType    = reflect.TypeOf(sql.NullTime{})
)
