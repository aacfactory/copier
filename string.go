package copier

import (
	"database/sql/driver"
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

func copyString(dst reflect.Value, src reflect.Value) (err error) {
	if src.Type().AssignableTo(dst.Type()) {
		dst.Set(src)
		return
	}
	switch src.Kind() {
	case reflect.String:
		dst.Set(src)
		break
	case reflect.Bool:
		v := src.Bool()
		if v {
			dst.SetString(trueString)
		} else {
			dst.SetString(falseString)
		}
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := src.Int()
		dst.SetString(strconv.FormatInt(v, 10))
		break
	case reflect.Float32, reflect.Float64:
		v := src.Float()
		dst.SetString(strconv.FormatFloat(v, 'f', 6, 64))
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := src.Uint()
		dst.SetString(strconv.FormatUint(v, 10))
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
			err = copyString(dst, reflect.ValueOf(value))
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
			dst.SetString(unsafe.String(unsafe.SliceData(p), len(p)))
			return
		}
		// time
		if src.Type().ConvertibleTo(timeType) {
			v := src.Interface().(time.Time)
			dst.SetString(v.Format(time.RFC3339))
			return
		}
		err = fmt.Errorf("copier: string not support %s source type", src.Type().String())
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
			dst.SetString(unsafe.String(unsafe.SliceData(p), len(p)))
			return
		}
		err = copyString(dst, src.Elem())
		break
	case reflect.Slice:
		if src.Type().Elem().Kind() == reflect.Uint8 {
			v := src.Interface().([]byte)
			dst.SetString(unsafe.String(unsafe.SliceData(v), len(v)))
			return
		}
		err = fmt.Errorf("copier: string not support %s source type", src.Type().String())
		break
	default:
		err = fmt.Errorf("copier: string not support %s source type", src.Type().String())
		return
	}
	return
}
