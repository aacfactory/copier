package copier

import (
	"database/sql/driver"
	"encoding"
	"fmt"
	"reflect"
	"unsafe"
)

func copyText(dst reflect.Value, src reflect.Value) (err error) {
	if src.Type().AssignableTo(dst.Type()) {
		dst.Set(src)
		return
	}
	switch src.Kind() {
	case reflect.String:
		v := src.String()
		if v == "" {
			return
		}
		p := unsafe.Slice(unsafe.StringData(v), len(v))
		method := dst.MethodByName("UnmarshalText")
		method.Call([]reflect.Value{reflect.ValueOf(p)})
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
			err = copyText(dst, reflect.ValueOf(value))
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
			err = copyText(dst, reflect.ValueOf(p))
			return
		}
		err = fmt.Errorf("copier: text not support %s source type", src.Type().String())
		break
	case reflect.Ptr:
		err = copyText(dst, src.Elem())
		break
	case reflect.Slice:
		if src.Type().Elem().Kind() == reflect.Uint8 {
			p := src.Interface().([]byte)
			method := dst.MethodByName("UnmarshalText")
			method.Call([]reflect.Value{reflect.ValueOf(p)})
			return
		}
		err = fmt.Errorf("copier: text not support %s source type", src.Type().String())
		break
	default:
		err = fmt.Errorf("copier: text not support %s source type", src.Type().String())
		return
	}
	return
}
