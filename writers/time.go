package writers

import (
	"database/sql/driver"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"time"
	"unsafe"
)

func NewTimeWriter(typ reflect2.Type) *TimeWriter {
	return &TimeWriter{
		typ: typ,
	}
}

type TimeWriter struct {
	typ reflect2.Type
}

func (w *TimeWriter) Name() string {
	return w.typ.String()
}

func (w *TimeWriter) Type() reflect2.Type {
	return w.typ
}

func (w *TimeWriter) Write(dst any, src any) (err error) {
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
	if srcType.Type1().ConvertibleTo(w.typ.Type1()) {
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}

	switch srcType.Kind() {
	case reflect.String:
		s := *(*string)(srcPtr)
		n, nErr := time.Parse(time.RFC3339, s)
		if nErr != nil {
			err = fmt.Errorf("copier: time writer can not support %s source type, src value is not RFC3339 format string", srcType.String())
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
	case reflect.Struct:
		// time
		if IsTime(srcType) {
			w.typ.UnsafeSet(dstPtr, srcPtr)
			break
		}
		// sql
		if valuer, ok := src.(driver.Valuer); ok {
			value, valueErr := valuer.Value()
			if valueErr != nil {
				err = valueErr
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
		err = fmt.Errorf("copier: time writer can not support %s source type", srcType.String())
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
		err = fmt.Errorf("copier: time writer can not support %s source type", srcType.String())
		return
	}
	return
}

func IsTime(typ reflect2.Type) bool {
	return timeType.RType() == typ.RType() || typ.Type1().ConvertibleTo(timeType.Type1())
}

func TimeToString(ptr unsafe.Pointer) unsafe.Pointer {
	v := new(time.Time)
	timeType.UnsafeSet(reflect2.PtrOf(v), ptr)
	s := v.Format(time.RFC3339)
	return reflect2.PtrOf(s)
}

func TimeToInt(ptr unsafe.Pointer) unsafe.Pointer {
	v := new(time.Time)
	timeType.UnsafeSet(reflect2.PtrOf(v), ptr)
	s := v.UnixMilli()
	return reflect2.PtrOf(s)
}
