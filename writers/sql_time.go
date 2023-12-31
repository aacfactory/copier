package writers

import (
	"database/sql/driver"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"time"
)

type SQLNullTimeWriter struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (w *SQLNullTimeWriter) Name() string {
	return w.typ.String()
}

func (w *SQLNullTimeWriter) Type() reflect2.Type {
	return w.typ
}

func (w *SQLNullTimeWriter) ValueType() reflect2.Type {
	return w.valueType.Type()
}

func (w *SQLNullTimeWriter) Write(dst any, src any) (err error) {
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
	if srcType.Type1().ConvertibleTo(w.valueType.Type().Type1()) {
		w.valueType.UnsafeSet(dstPtr, srcPtr)
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		return
	}

	switch srcType.Kind() {
	case reflect.String:
		s := *(*string)(srcPtr)
		v, parseErr := time.Parse(time.RFC3339, s)
		if parseErr != nil {
			err = parseErr
			return
		}
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(v))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n := *(*int64)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(time.UnixMilli(n)))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u := *(*uint64)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(time.UnixMilli(int64(u))))
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
		// convertable
		if convertible, ok := src.(Convertible); ok {
			value := convertible.Convert()
			if value == nil {
				return
			}
			err = w.Write(dst, reflect2.TypeOf(value).PackEFace(reflect2.PtrOf(value)))
			return
		}
		err = fmt.Errorf("copier: sql null time writer can not support %s type reader", srcType.String())
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
		err = fmt.Errorf("copier: sql null time writer can not support %s type reader", srcType.String())
		return
	}
	return
}
