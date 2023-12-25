package writers

import (
	"database/sql"
	"encoding"
	"fmt"
	"github.com/aacfactory/copier/internal/commons"
	"github.com/modern-go/reflect2"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

func NewStringWriter() Writer {
	return &StringWriter{
		typ: reflect2.TypeOf(""),
	}
}

type StringWriter struct {
	typ reflect2.Type
}

func (w *StringWriter) Write(dstPtr unsafe.Pointer, reader Reader) (err error) {
	srcType := reader.Type()
	srcPtr := reader.Read()
	switch srcType.Kind() {
	case reflect.String:
		w.typ.UnsafeSet(dstPtr, srcPtr)
		break
	case reflect.Bool:
		b := *(*bool)(srcPtr)
		s := strconv.FormatBool(b)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(s))
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n := *(*int64)(srcPtr)
		s := strconv.FormatInt(n, 10)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(s))
		break
	case reflect.Float32, reflect.Float64:
		f := *(*float64)(srcPtr)
		s := strconv.FormatFloat(f, 'f', 6, 64)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(s))
		break
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u := *(*uint64)(srcPtr)
		s := strconv.FormatUint(u, 10)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(s))
		break
	case reflect.Uint8:
		u := *(*uint8)(srcPtr)
		s := unsafe.String(unsafe.SliceData([]byte{u}), 1)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(s))
		break
	default:
		// sql: null string
		if srcType.Type1().ConvertibleTo(sqlNullStringType.Type1()) {
			nv := new(sql.NullString)
			nptr := reflect2.PtrOf(nv)
			sqlNullStringType.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(nv.String))
			}
			break
		}
		// sql: null bool
		if srcType.Type1().ConvertibleTo(sqlNullBoolType.Type1()) {
			nv := new(sql.NullBool)
			nptr := reflect2.PtrOf(nv)
			sqlNullBoolType.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				s := strconv.FormatBool(nv.Bool)
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(s))
			}
			break
		}
		// sql: null int16
		if srcType.Type1().ConvertibleTo(sqlNullInt16Type.Type1()) {
			nv := new(sql.NullInt16)
			nptr := reflect2.PtrOf(nv)
			sqlNullInt16Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				s := strconv.FormatInt(int64(nv.Int16), 10)
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(s))
			}
			break
		}
		// sql: null int32
		if srcType.Type1().ConvertibleTo(sqlNullInt32Type.Type1()) {
			nv := new(sql.NullInt32)
			nptr := reflect2.PtrOf(nv)
			sqlNullInt32Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				s := strconv.FormatInt(int64(nv.Int32), 10)
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(s))
			}
			break
		}
		// sql: null int64
		if srcType.Type1().ConvertibleTo(sqlNullInt64Type.Type1()) {
			nv := new(sql.NullInt64)
			nptr := reflect2.PtrOf(nv)
			sqlNullInt64Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				s := strconv.FormatInt(nv.Int64, 10)
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(s))
			}
			break
		}
		// sql: null float
		if srcType.Type1().ConvertibleTo(sqlNullFloat64Type.Type1()) {
			nv := new(sql.NullFloat64)
			nptr := reflect2.PtrOf(nv)
			sqlNullFloat64Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				s := strconv.FormatFloat(nv.Float64, 'f', 6, 64)
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(s))
			}
			break
		}
		// sql: null byte
		if srcType.Type1().ConvertibleTo(sqlNullByteType.Type1()) {
			nv := new(sql.NullByte)
			nptr := reflect2.PtrOf(nv)
			sqlNullByteType.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				s := unsafe.String(unsafe.SliceData([]byte{nv.Byte}), 1)
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(s))
			}
			break
		}
		// sql: null time
		if srcType.Type1().ConvertibleTo(sqlNullTimeType.Type1()) {
			nv := new(sql.NullTime)
			nptr := reflect2.PtrOf(nv)
			sqlNullTimeType.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				s := nv.Time.Format(time.RFC3339)
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(s))
			}
			break
		}
		// time
		if srcType.Type1().ConvertibleTo(timeType.Type1()) {
			nv := new(time.Time)
			nptr := reflect2.PtrOf(nv)
			timeType.UnsafeSet(nptr, srcPtr)
			s := nv.Format(time.RFC3339)
			w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(s))
			break
		}
		// bytes
		if srcType.AssignableTo(bytesType) {
			p := *(*[]byte)(srcPtr)
			s := unsafe.String(unsafe.SliceData(p), len(p))
			w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(s))
			break
		}
		// text
		if srcType.Implements(textMarshalerType) {
			src := srcType.UnsafeIndirect(srcPtr)
			if srcType.IsNullable() && reflect2.IsNil(src) {
				break
			}
			text := (src).(encoding.TextMarshaler)
			p, encodeErr := text.MarshalText()
			if encodeErr != nil {
				err = fmt.Errorf("copier: string writer can not support %s type reader, %v", srcType.String(), encodeErr)
				return
			}
			s := unsafe.String(unsafe.SliceData(p), len(p))
			w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(s))
			break
		}
		err = fmt.Errorf("copier: string writer can not support %s type reader", srcType.String())
		return
	}
	return
}

func NewNullStringWriter(typ reflect2.Type) Writer {
	if typ.Kind() == reflect.Ptr {
		typ = typ.(reflect2.PtrType).Elem()
	}
	descriptor := commons.DescribeStruct("", typ)
	var valueType reflect2.StructField
	var validType reflect2.StructField
	for _, field := range descriptor.Fields {
		for _, f := range field.Field {
			if f.Name() == "String" && f.Type().Kind() == reflect.String {
				valueType = f.Type().(reflect2.StructField)
				continue
			}
			if f.Name() == "Valid" && f.Type().Kind() == reflect.Bool {
				validType = f.Type().(reflect2.StructField)
				continue
			}
		}
	}
	return &NullStringWriter{
		typ:       typ,
		valueType: valueType,
		validType: validType,
	}
}

type NullStringWriter struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (w *NullStringWriter) Write(dstPtr unsafe.Pointer, reader Reader) (err error) {
	srcType := reader.Type()
	srcPtr := reader.Read()
	if w.typ.RType() == srcType.RType() {
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}
	switch srcType.Kind() {
	case reflect.String:
		value := *(*string)(srcPtr)
		if value == "" {
			return
		}
		w.valueType.UnsafeSet(dstPtr, srcPtr)
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Bool:
		value := *(*bool)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(strconv.FormatBool(value)))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value := *(*int64)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(strconv.FormatInt(value, 10)))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Float32, reflect.Float64:
		value := *(*float64)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(strconv.FormatFloat(value, 'f', 6, 64)))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value := *(*uint64)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(strconv.FormatUint(value, 10)))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Uint8:
		value := *(*uint8)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(unsafe.String(unsafe.SliceData([]byte{value}), 1)))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	default:
		// sql: null string
		if srcType.Type1().ConvertibleTo(sqlNullStringType.Type1()) {
			w.typ.UnsafeSet(dstPtr, srcPtr)
			break
		}
		// sql: null bool
		if srcType.Type1().ConvertibleTo(sqlNullBoolType.Type1()) {
			nv := new(sql.NullBool)
			nptr := reflect2.PtrOf(nv)
			sqlNullBoolType.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(strconv.FormatBool(nv.Bool)))
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
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(strconv.FormatInt(int64(nv.Int16), 10)))
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
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(strconv.FormatInt(int64(nv.Int32), 10)))
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
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(strconv.FormatInt(nv.Int64, 10)))
				w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		// sql: null float
		if srcType.Type1().ConvertibleTo(sqlNullFloat64Type.Type1()) {
			nv := new(sql.NullFloat64)
			nptr := reflect2.PtrOf(nv)
			sqlNullFloat64Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(strconv.FormatFloat(nv.Float64, 'f', 6, 64)))
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
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(unsafe.String(unsafe.SliceData([]byte{nv.Byte}), 1)))
				w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		// sql: null time
		if srcType.Type1().ConvertibleTo(sqlNullTimeType.Type1()) {
			nv := new(sql.NullTime)
			nptr := reflect2.PtrOf(nv)
			sqlNullTimeType.UnsafeSet(nptr, srcPtr)
			if nv.Valid && !nv.Time.IsZero() {
				s := nv.Time.Format(time.RFC3339)
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(s))
				w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		// time
		if srcType.Type1().ConvertibleTo(timeType.Type1()) {
			nv := new(time.Time)
			nptr := reflect2.PtrOf(nv)
			timeType.UnsafeSet(nptr, srcPtr)
			if !nv.IsZero() {
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(nv.Format(time.RFC3339)))
				w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		// bytes
		if srcType.AssignableTo(bytesType) {
			p := *(*[]byte)(srcPtr)
			if pLen := len(p); pLen > 0 {
				s := unsafe.String(unsafe.SliceData(p), pLen)
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(s))
				w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		// text
		if srcType.Implements(textMarshalerType) {
			src := srcType.UnsafeIndirect(srcPtr)
			if srcType.IsNullable() && reflect2.IsNil(src) {
				break
			}
			text := (src).(encoding.TextMarshaler)
			p, encodeErr := text.MarshalText()
			if encodeErr != nil {
				err = fmt.Errorf("copier: string writer can not support %s type reader, %v", srcType.String(), encodeErr)
				return
			}
			if pLen := len(p); pLen > 0 {
				s := unsafe.String(unsafe.SliceData(p), pLen)
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(s))
				w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		err = fmt.Errorf("copier: string writer can not support %s type reader", srcType.String())
		return
	}
	return
}
