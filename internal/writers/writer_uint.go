package writers

import (
	"database/sql"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

func NewUintWriter() Writer {
	return &IntWriter{
		typ: reflect2.TypeOf(uint64(0)),
	}
}

type UintWriter struct {
	typ reflect2.Type
}

func (w *UintWriter) Write(dstPtr unsafe.Pointer, reader Reader) (err error) {
	srcType := reader.Type()
	srcPtr := reader.Read()
	switch srcType.Kind() {
	case reflect.String:
		s := *(*string)(srcPtr)
		n, nErr := strconv.ParseUint(s, 10, 64)
		if nErr != nil {
			err = fmt.Errorf("copier: uint writer can not support %s type reader, src value is not uint format string", srcType.String())
			return
		}
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(n))
		break
	case reflect.Bool:
		b := *(*bool)(srcPtr)
		if b {
			w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(uint64(1)))
		}
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f := *(*int64)(srcPtr)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(uint64(f)))
		break
	case reflect.Float32, reflect.Float64:
		f := *(*float64)(srcPtr)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(uint64(f)))
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		w.typ.UnsafeSet(dstPtr, srcPtr)
		break
	default:
		// sql: null string
		if srcType.Type1().ConvertibleTo(sqlNullStringType.Type1()) {
			nv := new(sql.NullString)
			nptr := reflect2.PtrOf(nv)
			sqlNullStringType.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				n, nErr := strconv.ParseUint(nv.String, 10, 64)
				if nErr != nil {
					err = fmt.Errorf("copier: uint writer can not support %s type reader, src value is not uint format string", srcType.String())
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
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(uint64(1)))
			}
			break
		}
		// sql: null int16
		if srcType.Type1().ConvertibleTo(sqlNullInt16Type.Type1()) {
			nv := new(sql.NullInt16)
			nptr := reflect2.PtrOf(nv)
			sqlNullInt16Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(uint64(nv.Int16)))
			}
			break
		}
		// sql: null int32
		if srcType.Type1().ConvertibleTo(sqlNullInt32Type.Type1()) {
			nv := new(sql.NullInt32)
			nptr := reflect2.PtrOf(nv)
			sqlNullInt32Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(uint64(nv.Int32)))
			}
			break
		}
		// sql: null int64
		if srcType.Type1().ConvertibleTo(sqlNullInt64Type.Type1()) {
			nv := new(sql.NullInt64)
			nptr := reflect2.PtrOf(nv)
			sqlNullInt64Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(uint64(nv.Int64)))
			}
			break
		}
		// sql: null float
		if srcType.Type1().ConvertibleTo(sqlNullFloat64Type.Type1()) {
			nv := new(sql.NullFloat64)
			nptr := reflect2.PtrOf(nv)
			sqlNullFloat64Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(uint64(nv.Float64)))
			}
			break
		}
		// sql: null time
		if srcType.Type1().ConvertibleTo(sqlNullTimeType.Type1()) {
			nv := new(sql.NullTime)
			nptr := reflect2.PtrOf(nv)
			sqlNullTimeType.UnsafeSet(nptr, srcPtr)
			if nv.Valid && !nv.Time.IsZero() {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(uint64(nv.Time.UnixMilli())))
			}
			break
		}
		// time
		if srcType.Type1().ConvertibleTo(timeType.Type1()) {
			nv := new(time.Time)
			nptr := reflect2.PtrOf(nv)
			timeType.UnsafeSet(nptr, srcPtr)
			if !nv.IsZero() {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(uint64(nv.UnixMilli())))
			}
			break
		}
		err = fmt.Errorf("copier: uint writer can not support %s type reader", srcType.String())
		return
	}
	return
}
