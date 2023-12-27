package writers

import (
	"database/sql/driver"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"strconv"
	"unsafe"
)

func NewIntWriter() Writer {
	return &IntWriter{
		typ: intType,
	}
}

type IntWriter struct {
	typ reflect2.Type
}

func (w *IntWriter) Name() string {
	return "int"
}

func (w *IntWriter) Type() reflect2.Type {
	return w.typ
}

func (w *IntWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {
	switch srcType.Kind() {
	case reflect.String:
		s := *(*string)(srcPtr)
		n, nErr := strconv.ParseInt(s, 10, 64)
		if nErr != nil {
			err = fmt.Errorf("copier: int writer can not support %s source type, src value is not int format string", srcType.String())
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
		intType.UnsafeSet(dstPtr, srcPtr)
		break
	case reflect.Float32, reflect.Float64:
		f := *(*float64)(srcPtr)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(int64(f)))
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u := *(*uint64)(srcPtr)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(int64(u)))
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
	default:
		err = fmt.Errorf("copier: int writer can not support %s source type", srcType.String())
		return
	}
	return
}
