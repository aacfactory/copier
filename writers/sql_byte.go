package writers

import (
	"database/sql/driver"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"unsafe"
)

type SQLNullByteWriter struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (w *SQLNullByteWriter) Name() string {
	return w.typ.String()
}

func (w *SQLNullByteWriter) Type() reflect2.Type {
	return w.typ
}

func (w *SQLNullByteWriter) ValueType() reflect2.Type {
	return w.valueType.Type()
}

func (w *SQLNullByteWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {
	if w.typ.RType() == srcType.RType() {
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}
	switch srcType.Kind() {
	case reflect.Bool:
		b := *(*bool)(srcPtr)
		if b {
			w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf('t'))
		} else {
			w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf('f'))
		}
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n := *(*int64)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(byte(n)))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Uint8:
		u := *(*uint8)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(u))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.String:
		s := *(*string)(srcPtr)
		if s != "" {
			w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(s[0]))
		}
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Struct, reflect.Ptr:
		// sql
		if IsSQLValue(srcType) {
			valuer, isValuer := srcType.PackEFace(srcPtr).(driver.Valuer)
			if !isValuer {
				err = fmt.Errorf("copier: sql null byte writer can not support %s source type", srcType.String())
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
		err = fmt.Errorf("copier: sql null byte writer can not support %s source type", srcType.String())
		break
	}
	return
}
