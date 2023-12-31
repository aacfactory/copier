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

func NewStringWriter() Writer {
	return &StringWriter{
		typ: stringType,
	}
}

type StringWriter struct {
	typ reflect2.Type
}

func (w *StringWriter) Name() string {
	return "string"
}

func (w *StringWriter) Type() reflect2.Type {
	return w.typ
}

func (w *StringWriter) Write(dst any, src any) (err error) {
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

	switch srcType.Kind() {
	case reflect.String:
		w.typ.UnsafeSet(dstPtr, srcPtr)
		break
	case reflect.Bool:
		v := *(*bool)(srcPtr)
		s := falseString
		if v {
			s = trueString
		}
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(s))
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := *(*int64)(srcPtr)
		s := strconv.FormatInt(v, 10)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(s))
		break
	case reflect.Float32, reflect.Float64:
		v := *(*float64)(srcPtr)
		s := strconv.FormatFloat(v, 'f', 6, 64)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(s))
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := *(*uint64)(srcPtr)
		s := strconv.FormatUint(v, 10)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(s))
		break
	case reflect.Struct:
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
		// time
		if IsTime(srcType) {
			w.typ.UnsafeSet(dstPtr, TimeToString(srcPtr))
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
			w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(&s))
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
		err = fmt.Errorf("copier: string writer can not support %s source type", srcType.String())
		break
	case reflect.Ptr:
		// text
		if value, ok := src.(encoding.TextMarshaler); ok {
			p, encodeErr := value.MarshalText()
			if encodeErr != nil {
				err = encodeErr
				return
			}
			s := unsafe.String(unsafe.SliceData(p), len(p))
			w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(&s))
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
	case reflect.Slice:
		// bytes
		if srcType.(reflect2.SliceType).Elem().Kind() == reflect.Uint8 {
			p := *(*[]byte)(srcPtr)
			s := unsafe.String(unsafe.SliceData(p), len(p))
			w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(&s))
			return
		}
		err = fmt.Errorf("copier: string writer can not support %s source type", srcType.String())
		break
	default:
		err = fmt.Errorf("copier: string writer can not support %s source type", srcType.String())
		return
	}
	return
}
