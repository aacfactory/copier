package writers

import (
	"database/sql/driver"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"strconv"
	"unsafe"
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

func (w *BytesWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {
	// convertable
	if IsConvertible(srcType) {
		srcPtr, srcType = convert(srcPtr, srcType)
	}
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
	case reflect.Struct, reflect.Ptr:
		// time
		if IsTime(srcType) {
			w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(TimeToBytes(srcPtr)))
			break
		}
		// sql
		if IsSQLValue(srcType) {
			valuer, isValuer := srcType.PackEFace(srcPtr).(driver.Valuer)
			if !isValuer {
				err = fmt.Errorf("copier: bytes writer can not support %s source type", srcType.String())
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
			ptr, ptrErr := TextToBytes(srcPtr, srcType)
			if ptrErr != nil {
				err = ptrErr
				return
			}
			w.typ.UnsafeSet(dstPtr, ptr)
			break
		}
		err = fmt.Errorf("copier: bytes writer can not support %s source type", srcType.String())
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
