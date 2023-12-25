package descriptors

import (
	"database/sql"
	"encoding"
	"github.com/modern-go/reflect2"
	"math/big"
	"time"
)

var (
	timeType            = reflect2.TypeOf(time.Time{})
	bytesType           = reflect2.TypeOf([]byte{})
	bigIntType          = reflect2.TypeOf(big.NewInt(0))
	bigFloatType        = reflect2.TypeOf(big.NewFloat(0))
	bigRatType          = reflect2.TypeOf(big.NewRat(0, 1))
	sqlNullStringType   = reflect2.TypeOf(sql.NullString{})
	sqlNullBoolType     = reflect2.TypeOf(sql.NullBool{})
	sqlNullByteType     = reflect2.TypeOf(sql.NullByte{})
	sqlNullInt16Type    = reflect2.TypeOf(sql.NullInt16{})
	sqlNullInt32Type    = reflect2.TypeOf(sql.NullInt32{})
	sqlNullInt64Type    = reflect2.TypeOf(sql.NullInt64{})
	sqlNullFloat64Type  = reflect2.TypeOf(sql.NullFloat64{})
	sqlNullTimeType     = reflect2.TypeOf(sql.NullTime{})
	sqlScannerType      = reflect2.TypeOfPtr((*sql.Scanner)(nil)).Elem()
	textMarshalerType   = reflect2.TypeOfPtr((*encoding.TextMarshaler)(nil)).Elem()
	textUnmarshalerType = reflect2.TypeOfPtr((*encoding.TextUnmarshaler)(nil)).Elem()
)
