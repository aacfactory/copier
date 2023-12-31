package writers

import (
	"database/sql/driver"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"strconv"
	"time"
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

func (w *IntWriter) Write(dst any, src any) (err error) {
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
		if value, ok := src.(time.Time); ok {
			w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(value.UnixMilli()))
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
		err = fmt.Errorf("copier: int writer can not support %s source type", srcType.String())
		break
	case reflect.Ptr:
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
		err = fmt.Errorf("copier: int writer can not support %s source type", srcType.String())
		return
	}
	return
}
