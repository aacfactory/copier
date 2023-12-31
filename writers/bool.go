package writers

import (
	"database/sql/driver"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"strings"
)

const (
	trueString  = "true"
	falseString = "false"
)

func NewBoolWriter() Writer {
	return &BoolWriter{
		typ: boolType,
	}
}

type BoolWriter struct {
	typ reflect.Type
}

func (w *BoolWriter) Name() string {
	return "bool"
}

func (w *BoolWriter) Type() reflect2.Type {
	return w.typ
}

func (w *BoolWriter) Write(dst any, src any) (err error) {
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
	case reflect.Bool:
		w.typ.UnsafeSet(dstPtr, srcPtr)
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n := *(*int64)(srcPtr)
		if n > 0 {
			w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		}
		break
	case reflect.Uint8:
		u := *(*uint8)(srcPtr)
		if u == 'T' || u == 't' {
			w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		}
		break
	case reflect.String:
		s := *(*string)(srcPtr)
		if strings.ToLower(s) == trueString {
			w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		}
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
		// convertable
		if convertible, ok := src.(Convertible); ok {
			value := convertible.Convert()
			if value == nil {
				return
			}
			err = w.Write(dst, reflect2.TypeOf(value).PackEFace(reflect2.PtrOf(value)))
			return
		}
		err = fmt.Errorf("copier: bool writer can not support %s source type", srcType.String())
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
		err = fmt.Errorf("copier: bool writer can not support %s source type", srcType.String())
		break
	}
	return
}
