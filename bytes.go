package copier

import (
	"database/sql/driver"
	"encoding"
	"fmt"
	"reflect"
	"unsafe"
)

func copyBytes(dst reflect.Value, src reflect.Value) (err error) {
	if src.Type().AssignableTo(dst.Type()) {
		dst.Set(src)
		return
	}
	switch src.Kind() {
	case reflect.String:
		v := src.String()
		dst.SetBytes(unsafe.Slice(unsafe.StringData(v), len(v)))
		break
	case reflect.Struct:
		// text
		if src.Type().Implements(textMarshalerType) {
			v := src.Interface().(encoding.TextMarshaler)
			p, encodeErr := v.MarshalText()
			if encodeErr != nil {
				err = encodeErr
				return
			}
			dst.SetBytes(p)
			return
		}
		// sql
		if src.Type().Implements(sqlValuerType) {
			v := src.Interface().(driver.Valuer)
			value, valueErr := v.Value()
			if valueErr != nil {
				err = valueErr
				return
			}
			if value == nil {
				return
			}
			err = copyBytes(dst, reflect.ValueOf(value))
			return
		}
		err = fmt.Errorf("copier: bytes can not support %s source type", src.Type().String())
		break
	case reflect.Ptr:
		// text
		if src.Type().Implements(textMarshalerType) {
			v := src.Interface().(encoding.TextMarshaler)
			p, encodeErr := v.MarshalText()
			if encodeErr != nil {
				err = encodeErr
				return
			}
			dst.SetBytes(p)
			return
		}
		err = copyBytes(dst, src.Elem())
		break
	case reflect.Slice:
		if src.Type().Elem().Kind() == reflect.Uint8 {
			dst.Set(src)
			break
		}
		err = fmt.Errorf("copier: bytes can not support %s source type", src.Type().String())
		break
	default:
		err = fmt.Errorf("copier: bytes not support %s source type", src.Type().String())
		return
	}
	return
}
