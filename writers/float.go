package writers

import (
	"database/sql/driver"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"strconv"
	"unsafe"
)

func NewFloatWriter() Writer {
	return &FloatWriter{
		typ: floatType,
	}
}

type FloatWriter struct {
	typ reflect2.Type
}

func (w *FloatWriter) Name() string {
	return "float"
}

func (w *FloatWriter) Type() reflect2.Type {
	return w.typ
}

func (w *FloatWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {
	switch srcType.Kind() {
	case reflect.String:
		s := *(*string)(srcPtr)
		n, nErr := strconv.ParseFloat(s, 64)
		if nErr != nil {
			err = fmt.Errorf("copier: float writer can not support %s source type, src value is not float format string", srcType.String())
			return
		}
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(n))
		break
	case reflect.Bool:
		b := *(*bool)(srcPtr)
		if b {
			w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(float64(1)))
		}
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n := *(*int64)(srcPtr)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(float64(n)))
		break
	case reflect.Float32, reflect.Float64:
		w.typ.UnsafeSet(dstPtr, srcPtr)
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u := *(*uint64)(srcPtr)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(float64(u)))
		break
	case reflect.Struct, reflect.Ptr:
		// sql
		if IsSQLValue(srcType) {
			valuer, isValuer := srcType.PackEFace(srcPtr).(driver.Valuer)
			if !isValuer {
				err = fmt.Errorf("copier: float writer can not support %s source type", srcType.String())
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
		err = fmt.Errorf("copier: float writer can not support %s type reader", srcType.String())
		return
	}
	return
}
