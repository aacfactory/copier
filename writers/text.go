package writers

import (
	"database/sql/driver"
	"encoding"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
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

func (w *TextUnmarshalerWriter) Write(dst any, src any) (err error) {
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
	case reflect.Struct:
		// text
		if value, ok := src.(encoding.TextMarshaler); ok {
			p, encodeErr := value.MarshalText()
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
		err = fmt.Errorf("copier: text unmarshaler writer can not support %s source type", srcType.String())
		return
	case reflect.Ptr:
		// text
		if value, ok := src.(encoding.TextMarshaler); ok {
			p, encodeErr := value.MarshalText()
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
		err = fmt.Errorf("copier: text unmarshaler writer can not support %s source type", srcType.String())
		return
	}
	return
}

func IsText(typ reflect2.Type) bool {
	return typ.Implements(textMarshalerType)
}
