package writers

import (
	"database/sql/driver"
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

func (w *StringWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {
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
	case reflect.Struct, reflect.Ptr:
		// time
		if IsTime(srcType) {
			w.typ.UnsafeSet(dstPtr, TimeToString(srcPtr))
			break
		}
		// sql
		if IsSQLValue(srcType) {
			valuer, isValuer := srcType.PackEFace(srcPtr).(driver.Valuer)
			if !isValuer {
				err = fmt.Errorf("copier: string writer can not support %s source type", srcType.String())
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
		// text
		if IsText(srcType) {
			ptr, ptrErr := TextToString(srcPtr, srcType)
			if ptrErr != nil {
				err = ptrErr
				return
			}
			w.typ.UnsafeSet(dstPtr, ptr)
			break
		}
		err = fmt.Errorf("copier: string writer can not support %s source type", srcType.String())
		break
	case reflect.Slice:
		// bytes
		if srcType.(reflect2.SliceType).Elem().Kind() == reflect.Uint8 {
			p := *(*[]byte)(srcPtr)
			s := unsafe.String(unsafe.SliceData(p), len(p))
			w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(s))
			break
		}
		err = fmt.Errorf("copier: string writer can not support %s source type", srcType.String())
		break
	default:
		err = fmt.Errorf("copier: string writer can not support %s source type", srcType.String())
		return
	}
	return
}
