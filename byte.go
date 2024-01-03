package copier

import (
	"database/sql/driver"
	"encoding"
	"fmt"
	"reflect"
)

func copyByte(dst reflect.Value, src reflect.Value) (err error) {
	if src.Type().AssignableTo(dst.Type()) {
		dst.Set(src)
		return
	}
	switch src.Kind() {
	case reflect.String:
		v := src.String()
		if len(v) != 1 {
			err = fmt.Errorf("copier: byte can not support %s source type", src.Type().String())
			return
		}
		dst.Set(reflect.ValueOf(v[0]))
		break
	case reflect.Bool:
		v := src.Bool()

		if v {
			dst.Set(reflect.ValueOf('t'))
		} else {
			dst.Set(reflect.ValueOf('f'))
		}
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := src.Int()
		dst.Set(reflect.ValueOf(byte(v)))
		break
	case reflect.Uint8:
		v := src.Uint()
		dst.Set(reflect.ValueOf(byte(v)))
		break
	case reflect.Struct:
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
			err = copyByte(dst, reflect.ValueOf(value))
			return
		}
		// text
		if src.Type().Implements(textMarshalerType) {
			v := src.Interface().(encoding.TextMarshaler)
			p, encodeErr := v.MarshalText()
			if encodeErr != nil {
				err = encodeErr
				return
			}
			if len(p) != 1 {
				err = fmt.Errorf("copier: byte can not support %s source type", src.Type().String())
				return
			}
			dst.Set(reflect.ValueOf(p[0]))
			return
		}
		err = fmt.Errorf("copier: byte not support %s source type", src.Type().String())
		break
	case reflect.Slice:
		// bytes
		if src.Type().Elem().Kind() == reflect.Uint8 {
			v := src.Interface().([]byte)
			if len(v) != 1 {
				err = fmt.Errorf("copier: byte not support %s source type", src.Type().String())
				return
			}
			dst.Set(reflect.ValueOf(v[0]))
			return
		}
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
			if len(p) != 1 {
				err = fmt.Errorf("copier: byte can not support %s source type", src.Type().String())
				return
			}
			dst.Set(reflect.ValueOf(p[0]))
			return
		}
		err = copyByte(dst, src.Elem())
		break
	default:
		err = fmt.Errorf("copier: byte can not support %s source type", src.Type().String())
		break
	}
	return
}
