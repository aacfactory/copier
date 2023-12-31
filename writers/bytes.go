package writers

import (
	"database/sql/driver"
	"encoding"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"strconv"
)

func NewBytesWriter(typ reflect2.Type) *BytesWriter {
	return &BytesWriter{
		typ: typ,
	}
}

type BytesWriter struct {
	typ reflect2.Type
}

func (w *BytesWriter) Name() string {
	return w.typ.String()
}

func (w *BytesWriter) Type() reflect2.Type {
	return w.typ
}

func (w *BytesWriter) Write(dst any, src any) (err error) {
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
		v := *(*string)(srcPtr)
		p := reflect2.UnsafeCastString(v)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(p))
		break
	case reflect.Bool:
		v := *(*bool)(srcPtr)
		s := falseString
		if v {
			s = trueString
		}
		p := reflect2.UnsafeCastString(s)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(p))
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := *(*int64)(srcPtr)
		s := strconv.FormatInt(v, 10)
		p := reflect2.UnsafeCastString(s)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(p))
		break
	case reflect.Float32, reflect.Float64:
		v := *(*float64)(srcPtr)
		s := strconv.FormatFloat(v, 'f', 6, 64)
		p := reflect2.UnsafeCastString(s)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(p))
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := *(*uint64)(srcPtr)
		s := strconv.FormatUint(v, 10)
		p := reflect2.UnsafeCastString(s)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(p))
		break
	case reflect.Struct:
		// text
		if value, ok := src.(encoding.TextMarshaler); ok {
			p, encodeErr := value.MarshalText()
			if encodeErr != nil {
				err = encodeErr
				return
			}
			w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(&p))
			return
		}
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
		// convertable
		if convertible, ok := src.(Convertible); ok {
			value := convertible.Convert()
			if value == nil {
				return
			}
			err = w.Write(dst, reflect2.TypeOf(value).PackEFace(reflect2.PtrOf(value)))
			return
		}
		err = fmt.Errorf("copier: bytes writer can not support %s source type", srcType.String())
		break
	case reflect.Ptr:
		// text
		if value, ok := src.(encoding.TextMarshaler); ok {
			p, encodeErr := value.MarshalText()
			if encodeErr != nil {
				err = encodeErr
				return
			}
			w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(&p))
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
			w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(p))
			break
		}
		err = fmt.Errorf("copier: bytes writer can not support %s source type", srcType.String())
		break
	default:
		err = fmt.Errorf("copier: bytes writer can not support %s source type", srcType.String())
		return
	}
	return
}
