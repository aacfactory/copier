package writers

import (
	"database/sql"
	"fmt"
	"github.com/aacfactory/copier/internal/commons"
	"github.com/modern-go/reflect2"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

func NewIntWriter() Writer {
	return &IntWriter{
		typ: reflect2.TypeOf(int64(0)),
	}
}

type IntWriter struct {
	typ reflect2.Type
}

func (w *IntWriter) Write(dstPtr unsafe.Pointer, reader Reader) (err error) {
	srcType := reader.Type()
	srcPtr := reader.Read()
	switch srcType.Kind() {
	case reflect.String:
		s := *(*string)(srcPtr)
		n, nErr := strconv.ParseInt(s, 10, 64)
		if nErr != nil {
			err = fmt.Errorf("copier: int writer can not support %s type reader, src value is not int format string", srcType.String())
			return
		}
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(n))
		break
	case reflect.Bool:
		b := *(*bool)(srcPtr)
		if b {
			w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(1))
		}
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		w.typ.UnsafeSet(dstPtr, srcPtr)
		break
	case reflect.Float32, reflect.Float64:
		f := *(*float64)(srcPtr)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(int64(f)))
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u := *(*uint64)(srcPtr)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(int64(u)))
		break
	default:
		// sql: null string
		if srcType.Type1().ConvertibleTo(sqlNullStringType.Type1()) {
			nv := new(sql.NullString)
			nptr := reflect2.PtrOf(nv)
			sqlNullStringType.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				n, nErr := strconv.ParseInt(nv.String, 10, 64)
				if nErr != nil {
					err = fmt.Errorf("copier: int writer can not support %s type reader, src value is not int format string", srcType.String())
					return
				}
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(n))
			}
			break
		}
		// sql: null bool
		if srcType.Type1().ConvertibleTo(sqlNullBoolType.Type1()) {
			nv := new(sql.NullBool)
			nptr := reflect2.PtrOf(nv)
			sqlNullBoolType.UnsafeSet(nptr, srcPtr)
			if nv.Valid && nv.Bool {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(int64(1)))
			}
			break
		}
		// sql: null int16
		if srcType.Type1().ConvertibleTo(sqlNullInt16Type.Type1()) {
			nv := new(sql.NullInt16)
			nptr := reflect2.PtrOf(nv)
			sqlNullInt16Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(int64(nv.Int16)))
			}
			break
		}
		// sql: null int32
		if srcType.Type1().ConvertibleTo(sqlNullInt32Type.Type1()) {
			nv := new(sql.NullInt32)
			nptr := reflect2.PtrOf(nv)
			sqlNullInt32Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(int64(nv.Int32)))
			}
			break
		}
		// sql: null int64
		if srcType.Type1().ConvertibleTo(sqlNullInt64Type.Type1()) {
			nv := new(sql.NullInt64)
			nptr := reflect2.PtrOf(nv)
			sqlNullInt64Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(nv.Int64))
			}
			break
		}
		// sql: null float
		if srcType.Type1().ConvertibleTo(sqlNullFloat64Type.Type1()) {
			nv := new(sql.NullFloat64)
			nptr := reflect2.PtrOf(nv)
			sqlNullFloat64Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(int64(nv.Float64)))
			}
			break
		}
		// sql: null time
		if srcType.Type1().ConvertibleTo(sqlNullTimeType.Type1()) {
			nv := new(sql.NullTime)
			nptr := reflect2.PtrOf(nv)
			sqlNullTimeType.UnsafeSet(nptr, srcPtr)
			if nv.Valid && !nv.Time.IsZero() {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(nv.Time.UnixMilli()))
			}
			break
		}
		// time
		if srcType.Type1().ConvertibleTo(timeType.Type1()) {
			nv := new(time.Time)
			nptr := reflect2.PtrOf(nv)
			timeType.UnsafeSet(nptr, srcPtr)
			if !nv.IsZero() {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(nv.UnixMilli()))
			}
			break
		}
		err = fmt.Errorf("copier: int writer can not support %s type reader", srcType.String())
		return
	}
	return
}

func NewNullIntWriter(typ reflect2.Type) Writer {
	if typ.Kind() == reflect.Ptr {
		typ = typ.(reflect2.PtrType).Elem()
	}
	descriptor := commons.DescribeStruct("", typ)
	var valueType reflect2.StructField
	var validType reflect2.StructField
	for _, field := range descriptor.Fields {
		for _, f := range field.Field {
			if f.Name() == "Int16" && f.Type().Kind() == reflect.Int16 {
				valueType = f.Type().(reflect2.StructField)
				continue
			}
			if f.Name() == "Int32" && f.Type().Kind() == reflect.Int32 {
				valueType = f.Type().(reflect2.StructField)
				continue
			}
			if f.Name() == "Int64" && f.Type().Kind() == reflect.Int64 {
				valueType = f.Type().(reflect2.StructField)
				continue
			}
			if f.Name() == "Valid" && f.Type().Kind() == reflect.Bool {
				validType = f.Type().(reflect2.StructField)
				continue
			}
		}
	}
	return &NullIntWriter{
		typ:       typ,
		valueType: valueType,
		validType: validType,
	}
}

type NullIntWriter struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (w *NullIntWriter) Write(dstPtr unsafe.Pointer, reader Reader) (err error) {
	srcType := reader.Type()
	srcPtr := reader.Read()
	if w.typ.RType() == srcType.RType() {
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}
	switch srcType.Kind() {
	case reflect.String:
		s := *(*string)(srcPtr)
		n, nErr := strconv.ParseInt(s, 10, 64)
		if nErr != nil {
			err = fmt.Errorf("copier: null int writer can not support %s type reader, src value is not int format string", srcType.String())
			return
		}
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(n))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Bool:
		b := *(*bool)(srcPtr)
		if b {
			w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		}
		break
	case reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Int64:
		w.valueType.UnsafeSet(dstPtr, srcPtr)
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Float32, reflect.Float64:
		f := *(*float64)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(int64(f)))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u := *(*uint64)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(int64(u)))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	default:
		// sql: null string
		if srcType.Type1().ConvertibleTo(sqlNullStringType.Type1()) {
			nv := new(sql.NullString)
			nptr := reflect2.PtrOf(nv)
			sqlNullStringType.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				n, nErr := strconv.ParseInt(nv.String, 10, 64)
				if nErr != nil {
					err = fmt.Errorf("copier: null int writer can not support %s type reader, src value is not int format string", srcType.String())
					return
				}
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(n))
				w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		// sql: null bool
		if srcType.Type1().ConvertibleTo(sqlNullBoolType.Type1()) {
			nv := new(sql.NullBool)
			nptr := reflect2.PtrOf(nv)
			sqlNullBoolType.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				if nv.Bool {
					w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(1))
				}
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
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(int64(nv.Int16)))
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
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(int64(nv.Int32)))
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
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(nv.Int64))
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
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(int64(nv.Float64)))
				w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		// sql: null time
		if srcType.Type1().ConvertibleTo(sqlNullTimeType.Type1()) {
			nv := new(sql.NullTime)
			nptr := reflect2.PtrOf(nv)
			sqlNullTimeType.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(nv.Time.UnixMilli()))
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
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(nv.UnixMilli()))
				w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		err = fmt.Errorf("copier: null int writer can not support %s type reader", srcType.String())
		return
	}
	return
}
