package writers

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"reflect"
	"time"
)

var (
	stringType          = reflect.TypeOf("")
	boolType            = reflect.TypeOf(false)
	intType             = reflect.TypeOf(int64(0))
	floatType           = reflect.TypeOf(float64(0))
	uintType            = reflect.TypeOf(uint64(0))
	timeType            = reflect.TypeOf(time.Time{})
	bytesType           = reflect.TypeOf([]byte{})
	sqlNullStringType   = reflect.TypeOf(sql.NullString{})
	sqlNullBoolType     = reflect.TypeOf(sql.NullBool{})
	sqlNullByteType     = reflect.TypeOf(sql.NullByte{})
	sqlNullInt16Type    = reflect.TypeOf(sql.NullInt16{})
	sqlNullInt32Type    = reflect.TypeOf(sql.NullInt32{})
	sqlNullInt64Type    = reflect.TypeOf(sql.NullInt64{})
	sqlNullFloat64Type  = reflect.TypeOf(sql.NullFloat64{})
	sqlNullTimeType     = reflect.TypeOf(sql.NullTime{})
	sqlValuerType       = reflect.TypeOf((*driver.Valuer)(nil)).Elem()
	textMarshalerType   = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
	textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
	unsafeWriterType    = reflect.TypeOf((*GenericWriter)(nil)).Elem()
)
