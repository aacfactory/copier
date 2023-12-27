package writers

import (
	"database/sql/driver"
	"encoding"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"unsafe"
)

func NewTextUnmarshalerWriter(typ reflect2.Type) *TextUnmarshalerWriter {
	return &TextUnmarshalerWriter{
		typ: typ,
	}
}

type TextUnmarshalerWriter struct {
	typ reflect2.Type
}

func (w *TextUnmarshalerWriter) Name() string {
	return w.typ.String()
}

func (w *TextUnmarshalerWriter) Type() reflect2.Type {
	return w.typ
}

func (w *TextUnmarshalerWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {
	if w.typ.RType() == srcType.RType() {
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}
	switch srcType.Kind() {
	case reflect.String:
		s := *(*string)(srcPtr)
		p := reflect2.UnsafeCastString(s)
		if len(p) > 0 {
			unmarshaler := w.typ.PackEFace(dstPtr).(encoding.TextUnmarshaler)
			if err = unmarshaler.UnmarshalText(p); err != nil {
				return
			}
		}
		break
	case reflect.Slice:
		if srcType.(reflect2.SliceType).Elem().Kind() == reflect.Uint8 {
			p := *(*[]byte)(srcPtr)
			if len(p) > 0 {
				unmarshaler := w.typ.PackEFace(dstPtr).(encoding.TextUnmarshaler)
				if err = unmarshaler.UnmarshalText(p); err != nil {
					return
				}
			}
			break
		}
		err = fmt.Errorf("copier: text unmarshaler writer can not support %s source type", srcType.String())
		break
	case reflect.Struct, reflect.Ptr:
		// text
		if IsText(srcType) {
			marshaler := srcType.PackEFace(srcPtr).(encoding.TextMarshaler)
			p, encodeErr := marshaler.MarshalText()
			if encodeErr != nil {
				err = encodeErr
				return
			}
			if len(p) > 0 {
				unmarshaler := w.typ.PackEFace(dstPtr).(encoding.TextUnmarshaler)
				if err = unmarshaler.UnmarshalText(p); err != nil {
					return
				}
			}
			break
		}
		// sql
		if IsSQLValue(srcType) {
			valuer, isValuer := srcType.PackEFace(srcPtr).(driver.Valuer)
			if !isValuer {
				err = fmt.Errorf("copier: text unmarshaler writer can not support %s source type", srcType.String())
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
		err = fmt.Errorf("copier: text unmarshaler writer can not support %s source type", srcType.String())
		return
	}
	return
}

func IsText(typ reflect2.Type) bool {
	return typ.Implements(textMarshalerType)
}

func TextToString(ptr unsafe.Pointer, typ reflect2.Type) (unsafe.Pointer, error) {
	src := typ.UnsafeIndirect(ptr)
	if typ.IsNullable() && typ.UnsafeIsNil(ptr) {
		return reflect2.PtrOf(""), nil
	}
	text := (src).(encoding.TextMarshaler)
	p, encodeErr := text.MarshalText()
	if encodeErr != nil {
		return nil, encodeErr
	}
	s := unsafe.String(unsafe.SliceData(p), len(p))
	return reflect2.PtrOf(s), nil
}

func TextToBytes(ptr unsafe.Pointer, typ reflect2.Type) (unsafe.Pointer, error) {
	src := typ.UnsafeIndirect(ptr)
	if typ.IsNullable() && typ.UnsafeIsNil(ptr) {
		return reflect2.PtrOf([]byte{}), nil
	}
	text := (src).(encoding.TextMarshaler)
	p, encodeErr := text.MarshalText()
	if encodeErr != nil {
		return nil, encodeErr
	}
	return reflect2.PtrOf(p), nil
}
