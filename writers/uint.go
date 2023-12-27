package writers

import (
	"database/sql/driver"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"strconv"
	"unsafe"
)

func NewUintWriter() Writer {
	return &UintWriter{
		typ: uintType,
	}
}

type UintWriter struct {
	typ reflect2.Type
}

func (w *UintWriter) Name() string {
	return "uint"
}

func (w *UintWriter) Type() reflect2.Type {
	return w.typ
}

func (w *UintWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {
	switch srcType.Kind() {
	case reflect.String:
		s := *(*string)(srcPtr)
		n, nErr := strconv.ParseUint(s, 10, 64)
		if nErr != nil {
			err = fmt.Errorf("copier: uint writer can not support %s source type, src value is not uint format string", srcType.String())
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
	case reflect.Struct, reflect.Ptr:
		// time
		if IsTime(srcType) {
			w.typ.UnsafeSet(dstPtr, TimeToInt(srcPtr))
			break
		}
		// sql
		if IsSQLValue(srcType) {
			valuer, isValuer := srcType.PackEFace(srcPtr).(driver.Valuer)
			if !isValuer {
				err = fmt.Errorf("copier: int writer can not support %s source type", srcType.String())
				return
			}
			value, valueErr := valuer.Value()
			if valueErr != nil {
				err = valueErr
				return
			}
			err = w.Write(dstPtr, reflect2.PtrOf(value), reflect2.TypeOf(value))
			return
		}
		break
	default:
		err = fmt.Errorf("copier: uint writer can not support %s source reader", srcType.String())
		return
	}
	return
}
