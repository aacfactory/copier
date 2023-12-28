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

func (w *TimeWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {
	if w.typ.RType() == srcType.RType() {
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}
	// convertible
	if w.typ.Type1().ConvertibleTo(srcType.Type1()) {
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}
	// convertable
	if IsConvertible(srcType) {
		srcPtr, srcType = convert(srcPtr, srcType)
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
	case reflect.Struct, reflect.Ptr:
		// time
		if IsTime(srcType) {
			w.typ.UnsafeSet(dstPtr, srcPtr)
			break
		}
		// sql
		if IsSQLValue(srcType) {
			valuer, isValuer := srcType.PackEFace(srcPtr).(driver.Valuer)
			if !isValuer {
				err = fmt.Errorf("copier: time writer can not support %s source type", srcType.String())
				return
			}
			value, valueErr := valuer.Value()
			if valueErr != nil {
				err = valueErr
				return
			}
			if reflect2.IsNil(value) {
				return
			}
			err = w.Write(dstPtr, reflect2.PtrOf(value), reflect2.TypeOf(value))
			return
		}
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

func TimeToBytes(ptr unsafe.Pointer) unsafe.Pointer {
	v := new(time.Time)
	timeType.UnsafeSet(reflect2.PtrOf(v), ptr)
	s := v.Format(time.RFC3339)
	return reflect2.PtrOf(reflect2.UnsafeCastString(s))
}

func TimeToInt(ptr unsafe.Pointer) unsafe.Pointer {
	v := new(time.Time)
	timeType.UnsafeSet(reflect2.PtrOf(v), ptr)
	s := v.UnixMilli()
	return reflect2.PtrOf(s)
}
