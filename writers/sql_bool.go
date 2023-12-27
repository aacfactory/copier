package writers

import (
	"database/sql/driver"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"strings"
	"unsafe"
)

type SQLNullBoolWriter struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (w *SQLNullBoolWriter) Name() string {
	return w.typ.String()
}

func (w *SQLNullBoolWriter) Type() reflect2.Type {
	return w.typ
}

func (w *SQLNullBoolWriter) ValueType() reflect2.Type {
	return w.valueType.Type()
}

func (w *SQLNullBoolWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {
	if w.typ.RType() == srcType.RType() {
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}
	switch srcType.Kind() {
	case reflect.Bool:
		w.valueType.UnsafeSet(dstPtr, srcPtr)
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n := *(*int64)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(n > 0))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Uint8:
		u := *(*uint8)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(u == 'T' || u == 't'))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.String:
		s := *(*string)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(strings.ToLower(s) == trueString))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Struct, reflect.Ptr:
		// sql
		if IsSQLValue(srcType) {
			valuer, isValuer := srcType.PackEFace(srcPtr).(driver.Valuer)
			if !isValuer {
				err = fmt.Errorf("copier: sql null bool writer can not support %s source type", srcType.String())
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
		break
	default:
		err = fmt.Errorf("copier: sql null bool writer can not support %s source type", srcType.String())
		break
	}
	return
}
