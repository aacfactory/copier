package writers

import (
	"database/sql"
	"fmt"
	"github.com/aacfactory/copier/internal/commons"
	"github.com/modern-go/reflect2"
	"reflect"
	"time"
	"unsafe"
)

func NewTimeWriter() Writer {
	return &TimeWriter{
		typ: reflect2.TypeOf(time.Time{}),
	}
}

type TimeWriter struct {
	typ reflect2.Type
}

func (w *TimeWriter) Write(dstPtr unsafe.Pointer, reader Reader) (err error) {
	srcType := reader.Type()
	srcPtr := reader.Read()
	if w.typ.RType() == srcType.RType() {
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}
	switch srcType.Kind() {
	case reflect.String:
		s := *(*string)(srcPtr)
		n, nErr := time.Parse(time.RFC3339, s)
		if nErr != nil {
			err = fmt.Errorf("copier: time writer can not support %s type reader, src value is not RFC3339 format string", srcType.String())
			return
		}
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(n))
		break
	case reflect.Int, reflect.Int64:
		n := *(*int64)(srcPtr)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(time.UnixMilli(n)))
		break
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		n := *(*uint64)(srcPtr)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(time.UnixMilli(int64(n))))
		break
	default:
		// sql: null string
		if srcType.Type1().ConvertibleTo(sqlNullStringType.Type1()) {
			nv := new(sql.NullString)
			nptr := reflect2.PtrOf(nv)
			sqlNullStringType.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				n, nErr := time.Parse(time.RFC3339, nv.String)
				if nErr != nil {
					err = fmt.Errorf("copier: time writer can not support %s type reader, src value is not RFC3339 format string", srcType.String())
					return
				}
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(n))
			}
			break
		}
		// sql: null int64
		if srcType.Type1().ConvertibleTo(sqlNullInt64Type.Type1()) {
			nv := new(sql.NullInt64)
			nptr := reflect2.PtrOf(nv)
			sqlNullInt64Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(time.UnixMilli(nv.Int64)))
			}
			break
		}
		// sql: null time
		if srcType.Type1().ConvertibleTo(sqlNullTimeType.Type1()) {
			nv := new(sql.NullTime)
			nptr := reflect2.PtrOf(nv)
			sqlNullTimeType.UnsafeSet(nptr, srcPtr)
			if nv.Valid && !nv.Time.IsZero() {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(nv.Time))
			}
			break
		}
		// time
		if srcType.Type1().ConvertibleTo(timeType.Type1()) {
			nv := new(time.Time)
			nptr := reflect2.PtrOf(nv)
			timeType.UnsafeSet(nptr, srcPtr)
			if !nv.IsZero() {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(nv))
			}
			break
		}
		err = fmt.Errorf("copier: time writer can not support %s type reader", srcType.String())
		return
	}
	return
}

func NewNullTimeWriter(typ reflect2.Type) Writer {
	if typ.Kind() == reflect.Ptr {
		typ = typ.(reflect2.PtrType).Elem()
	}
	descriptor := commons.DescribeStruct("", typ)
	var valueType reflect2.StructField
	var validType reflect2.StructField
	for _, field := range descriptor.Fields {
		for _, f := range field.Field {
			if f.Name() == "Time" && f.Type().Type1().ConvertibleTo(timeType.Type1()) {
				valueType = f.Type().(reflect2.StructField)
				continue
			}
			if f.Name() == "Valid" && f.Type().Kind() == reflect.Bool {
				validType = f.Type().(reflect2.StructField)
				continue
			}
		}
	}
	return &NullTimeWriter{
		typ:       typ,
		valueType: valueType,
		validType: validType,
	}
}

type NullTimeWriter struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (w *NullTimeWriter) Write(dstPtr unsafe.Pointer, reader Reader) (err error) {
	srcType := reader.Type()
	srcPtr := reader.Read()
	if w.typ.RType() == srcType.RType() {
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}
	switch srcType.Kind() {
	case reflect.String:
		s := *(*string)(srcPtr)
		n, nErr := time.Parse(time.RFC3339, s)
		if nErr != nil {
			err = fmt.Errorf("copier: null time writer can not support %s type reader, src value is not RFC3339 format string", srcType.String())
			return
		}
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(n))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Int, reflect.Int64:
		n := *(*int64)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(time.UnixMilli(n)))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		n := *(*uint64)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(time.UnixMilli(int64(n))))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	default:
		// sql: null string
		if srcType.Type1().ConvertibleTo(sqlNullStringType.Type1()) {
			nv := new(sql.NullString)
			nptr := reflect2.PtrOf(nv)
			sqlNullStringType.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				n, nErr := time.Parse(time.RFC3339, nv.String)
				if nErr != nil {
					err = fmt.Errorf("copier: null time writer can not support %s type reader, src value is not RFC3339 format string", srcType.String())
					return
				}
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(n))
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
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(time.UnixMilli(nv.Int64)))
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
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(nv.Time))
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
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(nv))
				w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		err = fmt.Errorf("copier: null time writer can not support %s type reader", srcType.String())
		return
	}
	return
}
