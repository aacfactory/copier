package writers

import (
	"database/sql/driver"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"strconv"
	"time"
)

type SQLNullIntWriter struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (w *SQLNullIntWriter) Name() string {
	return w.typ.String()
}

func (w *SQLNullIntWriter) Type() reflect2.Type {
	return w.typ
}

func (w *SQLNullIntWriter) ValueType() reflect2.Type {
	return w.valueType.Type()
}

func (w *SQLNullIntWriter) Write(dst any, src any) (err error) {
	if src == nil {
		return
	}
	srcType := reflect2.TypeOfPtr(src).Elem()
	srcPtr := reflect2.PtrOf(src)
	dstPtr := reflect2.PtrOf(dst)

	if w.typ.RType() == srcType.RType() {
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}
	if w.valueType.Type().RType() == srcType.RType() {
		w.valueType.UnsafeSet(dstPtr, srcPtr)
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		return
	}

	switch srcType.Kind() {
	case reflect.String:
		s := *(*string)(srcPtr)
		n, nErr := strconv.ParseInt(s, 10, 64)
		if nErr != nil {
			err = fmt.Errorf("copier: sql null int writer can not support %s source type, src value is not float format string", srcType.String())
			return
		}
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(n))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Bool:
		b := *(*bool)(srcPtr)
		if b {
			w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(1))
		}
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		w.valueType.UnsafeSet(dstPtr, srcPtr)
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Float32, reflect.Float64:
		n := *(*float64)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(int64(n)))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u := *(*uint64)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(int64(u)))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Struct:
		// sql
		if valuer, ok := src.(driver.Valuer); ok {
			value, valueErr := valuer.Value()
			if valueErr != nil {
				err = valueErr
				return
			}
			if value == nil {
				return
			}
			err = w.Write(dst, reflect2.TypeOf(value).PackEFace(reflect2.PtrOf(value)))
			return
		}
		// time
		if value, ok := src.(time.Time); ok {
			if !value.IsZero() {
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(value.UnixMilli()))
				w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			return
		}
		// convertable
		if convertible, ok := src.(Convertible); ok {
			value := convertible.Convert()
			if value == nil {
				return
			}
			err = w.Write(dst, reflect2.TypeOf(value).PackEFace(reflect2.PtrOf(value)))
			return
		}
		err = fmt.Errorf("copier: sql null int writer can not support %s type reader", srcType.String())
		return
	case reflect.Ptr:
		// convertable
		if convertible, ok := src.(Convertible); ok {
			value := convertible.Convert()
			if value == nil {
				return
			}
			err = w.Write(dst, reflect2.TypeOf(value).PackEFace(reflect2.PtrOf(value)))
			return
		}
		err = w.Write(dst, srcType.Indirect(src))
		break
	default:
		err = fmt.Errorf("copier: sql null int writer can not support %s type reader", srcType.String())
		return
	}
	return
}
