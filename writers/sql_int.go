package writers

import (
	"database/sql/driver"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"strconv"
	"unsafe"
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

func (w *SQLNullIntWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {
	if w.typ.RType() == srcType.RType() {
		w.typ.UnsafeSet(dstPtr, srcPtr)
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
	case reflect.Struct, reflect.Ptr:
		// sql
		if IsSQLValue(srcType) {
			valuer, isValuer := srcType.PackEFace(srcPtr).(driver.Valuer)
			if !isValuer {
				err = fmt.Errorf("copier: sql null int writer can not support %s source type", srcType.String())
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
		err = fmt.Errorf("copier: sql null int writer can not support %s type reader", srcType.String())
		return
	}
	return
}
