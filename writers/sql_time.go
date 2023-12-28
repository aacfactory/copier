package writers

import (
	"database/sql/driver"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"time"
	"unsafe"
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

func (w *SQLNullTimeWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {
	if w.typ.RType() == srcType.RType() {
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
	case reflect.Struct, reflect.Ptr:
		// sql
		if IsSQLValue(srcType) {
			valuer, isValuer := srcType.PackEFace(srcPtr).(driver.Valuer)
			if !isValuer {
				err = fmt.Errorf("copier: sql null time writer can not support %s source type", srcType.String())
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
		err = fmt.Errorf("copier: sql null time writer can not support %s type reader", srcType.String())
		return
	}
	return
}
