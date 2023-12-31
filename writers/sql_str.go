package writers

import (
	"database/sql/driver"
	"encoding"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"strconv"
	"unsafe"
)

type SQLNullStringWriter struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (w *SQLNullStringWriter) Name() string {
	return w.typ.String()
}

func (w *SQLNullStringWriter) Type() reflect2.Type {
	return w.typ
}

func (w *SQLNullStringWriter) ValueType() reflect2.Type {
	return w.valueType.Type()
}

func (w *SQLNullStringWriter) Write(dst any, src any) (err error) {
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

	switch srcType.Kind() {
	case reflect.String:
		w.valueType.UnsafeSet(dstPtr, srcPtr)
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Bool:
		b := *(*bool)(srcPtr)
		if b {
			w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(trueString))
		} else {
			w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(falseString))
		}
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n := *(*int64)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(strconv.FormatInt(n, 10)))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Float32, reflect.Float64:
		n := *(*float64)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(strconv.FormatFloat(n, 'f', 6, 64)))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u := *(*uint64)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(strconv.FormatUint(u, 10)))
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
		// text
		if value, ok := src.(encoding.TextMarshaler); ok {
			p, encodeErr := value.MarshalText()
			if encodeErr != nil {
				err = encodeErr
				return
			}
			s := unsafe.String(unsafe.SliceData(p), len(p))
			w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(&s))
			w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
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
		err = fmt.Errorf("copier: sql null string writer can not support %s type reader", srcType.String())
		return
	case reflect.Ptr:
		// text
		if value, ok := src.(encoding.TextMarshaler); ok {
			p, encodeErr := value.MarshalText()
			if encodeErr != nil {
				err = encodeErr
				return
			}
			s := unsafe.String(unsafe.SliceData(p), len(p))
			w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(&s))
			w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
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
		err = w.Write(dst, srcType.Indirect(src))
		break
	default:
		err = fmt.Errorf("copier: sql null string writer can not support %s type reader", srcType.String())
		return
	}
	return
}
