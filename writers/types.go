package writers

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"github.com/modern-go/reflect2"
	"time"
)

var (
	stringType          = reflect2.TypeOf("")
	boolType            = reflect2.TypeOf(false)
	intType             = reflect2.TypeOf(int64(0))
	floatType           = reflect2.TypeOf(float64(0))
	uintType            = reflect2.TypeOf(uint64(0))
	timeType            = reflect2.TypeOf(time.Time{})
	bytesType           = reflect2.TypeOf([]byte{})
	sqlNullStringType   = reflect2.TypeOf(sql.NullString{})
	sqlNullBoolType     = reflect2.TypeOf(sql.NullBool{})
	sqlNullByteType     = reflect2.TypeOf(sql.NullByte{})
	sqlNullInt16Type    = reflect2.TypeOf(sql.NullInt16{})
	sqlNullInt32Type    = reflect2.TypeOf(sql.NullInt32{})
	sqlNullInt64Type    = reflect2.TypeOf(sql.NullInt64{})
	sqlNullFloat64Type  = reflect2.TypeOf(sql.NullFloat64{})
	sqlNullTimeType     = reflect2.TypeOf(sql.NullTime{})
	sqlValuerType       = reflect2.TypeOfPtr((*driver.Valuer)(nil)).Elem()
	textMarshalerType   = reflect2.TypeOfPtr((*encoding.TextMarshaler)(nil)).Elem()
	textUnmarshalerType = reflect2.TypeOfPtr((*encoding.TextUnmarshaler)(nil)).Elem()
	unsafeWriterType    = reflect2.TypeOfPtr((*GenericWriter)(nil)).Elem()
	convertibleType     = reflect2.TypeOfPtr((*Convertible)(nil)).Elem()
)
