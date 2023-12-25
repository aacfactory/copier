package writers

import (
	"database/sql"
	"fmt"
	"github.com/aacfactory/copier/internal/commons"
	"github.com/modern-go/reflect2"
	"reflect"
	"strings"
	"unsafe"
)

func NewBoolWriter() Writer {
	return &BoolWriter{
		typ: reflect2.TypeOf(true),
	}
}

type BoolWriter struct {
	typ reflect2.Type
}

func (w *BoolWriter) Write(dstPtr unsafe.Pointer, reader Reader) (err error) {
	srcType := reader.Type()
	srcPtr := reader.Read()
	switch srcType.Kind() {
	case reflect.Bool:
		w.typ.UnsafeSet(dstPtr, srcPtr)
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n := *(*int64)(srcPtr)
		if n > 0 {
			w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		}
		break
	case reflect.Uint8:
		u := *(*uint8)(srcPtr)
		if u == 'T' || u == 't' {
			w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		}
		break
	case reflect.String:
		s := *(*string)(srcPtr)
		if strings.ToLower(s) == "true" {
			w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		}
		break
	default:
		// sql: null bool
		if srcType.Type1().ConvertibleTo(sqlNullBoolType.Type1()) {
			nv := new(sql.NullBool)
			nptr := reflect2.PtrOf(nv)
			sqlNullBoolType.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(nv.Bool))
			}
			break
		}
		// sql: null string
		if srcType.Type1().ConvertibleTo(sqlNullStringType.Type1()) {
			nv := new(sql.NullString)
			nptr := reflect2.PtrOf(nv)
			sqlNullStringType.UnsafeSet(nptr, srcPtr)
			if nv.Valid && strings.ToLower(nv.String) == "true" {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		// sql: null int16
		if srcType.Type1().ConvertibleTo(sqlNullInt16Type.Type1()) {
			nv := new(sql.NullInt16)
			nptr := reflect2.PtrOf(nv)
			sqlNullInt16Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid && nv.Int16 > 0 {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		// sql: null int32
		if srcType.Type1().ConvertibleTo(sqlNullInt32Type.Type1()) {
			nv := new(sql.NullInt32)
			nptr := reflect2.PtrOf(nv)
			sqlNullInt32Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid && nv.Int32 > 0 {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		// sql: null int64
		if srcType.Type1().ConvertibleTo(sqlNullInt64Type.Type1()) {
			nv := new(sql.NullInt64)
			nptr := reflect2.PtrOf(nv)
			sqlNullInt64Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid && nv.Int64 > 0 {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		// sql: null byte
		if srcType.Type1().ConvertibleTo(sqlNullByteType.Type1()) {
			nv := new(sql.NullByte)
			nptr := reflect2.PtrOf(nv)
			sqlNullByteType.UnsafeSet(nptr, srcPtr)
			if nv.Valid && (nv.Byte == 'T' || nv.Byte == 't') {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		err = fmt.Errorf("copier: bool writer can not support %s type reader", srcType.String())
		break
	}
	return
}

func NewNullBoolWriter(typ reflect2.Type) Writer {
	if typ.Kind() == reflect.Ptr {
		typ = typ.(reflect2.PtrType).Elem()
	}
	descriptor := commons.DescribeStruct("", typ)
	var valueType reflect2.StructField
	var validType reflect2.StructField
	for _, field := range descriptor.Fields {
		for _, f := range field.Field {
			if f.Name() == "Bool" && f.Type().Kind() == reflect.Bool {
				valueType = f.Type().(reflect2.StructField)
				continue
			}
			if f.Name() == "Valid" && f.Type().Kind() == reflect.Bool {
				validType = f.Type().(reflect2.StructField)
				continue
			}
		}
	}
	return &NullBoolWriter{
		typ:       typ,
		valueType: valueType,
		validType: validType,
	}
}

type NullBoolWriter struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (w *NullBoolWriter) Write(dstPtr unsafe.Pointer, reader Reader) (err error) {
	srcType := reader.Type()
	srcPtr := reader.Read()
	if w.typ.RType() == srcType.RType() {
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}
	switch srcType.Kind() {
	case reflect.Bool:
		value := *(*bool)(srcPtr)
		if value {
			w.valueType.UnsafeSet(dstPtr, srcPtr)
			w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			return
		}
		break
	case reflect.String:
		value := *(*string)(srcPtr)
		if strings.ToLower(value) == "true" {
			w.valueType.UnsafeSet(dstPtr, srcPtr)
			w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			return
		}
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value := *(*int64)(srcPtr)
		if value > 0 {
			w.valueType.UnsafeSet(dstPtr, srcPtr)
			w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			return
		}
		break
	case reflect.Uint8:
		value := *(*uint8)(srcPtr)
		if value == 'T' || value == 't' {
			w.valueType.UnsafeSet(dstPtr, srcPtr)
			w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			return
		}
		break
	default:
		// sql: null string
		if srcType.Type1().ConvertibleTo(sqlNullStringType.Type1()) {
			nv := new(sql.NullString)
			nptr := reflect2.PtrOf(nv)
			sqlNullStringType.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(nv.String))
				w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		// sql: null int16
		if srcType.Type1().ConvertibleTo(sqlNullInt16Type.Type1()) {
			nv := new(sql.NullInt16)
			nptr := reflect2.PtrOf(nv)
			sqlNullInt16Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(nv.Int16 > 0))
				w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		// sql: null int32
		if srcType.Type1().ConvertibleTo(sqlNullInt32Type.Type1()) {
			nv := new(sql.NullInt32)
			nptr := reflect2.PtrOf(nv)
			sqlNullInt32Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(nv.Int32 > 0))
				w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		// sql: null int64
		if srcType.Type1().ConvertibleTo(sqlNullInt64Type.Type1()) {
			nv := new(sql.NullInt64)
			nptr := reflect2.PtrOf(nv)
			sqlNullInt64Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(nv.Int64 > 0))
				w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		// sql: null byte
		if srcType.Type1().ConvertibleTo(sqlNullByteType.Type1()) {
			nv := new(sql.NullByte)
			nptr := reflect2.PtrOf(nv)
			sqlNullByteType.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(nv.Byte == 'T' || nv.Byte == 't'))
				w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		err = fmt.Errorf("copier: null bool writer can not support %s type reader", srcType.String())
		return
	}
	return
}
