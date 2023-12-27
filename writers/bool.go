package writers

import (
	"database/sql/driver"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"strings"
	"unsafe"
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
	typ reflect2.Type
}

func (w *BoolWriter) Name() string {
	return "bool"
}

func (w *BoolWriter) Type() reflect2.Type {
	return w.typ
}

func (w *BoolWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {
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
	case reflect.Struct, reflect.Ptr:
		// sql
		if IsSQLValue(srcType) {
			valuer, isValuer := srcType.PackEFace(srcPtr).(driver.Valuer)
			if !isValuer {
				err = fmt.Errorf("copier: bool writer can not support %s source type", srcType.String())
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
		err = fmt.Errorf("copier: bool writer can not support %s source type", srcType.String())
		break
	}
	return
}
