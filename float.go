package copier

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
)

func copyFloat(dst reflect.Value, src reflect.Value) (err error) {
	if src.Type().AssignableTo(dst.Type()) {
		dst.Set(src)
		return
	}
	switch src.Kind() {
	case reflect.String:
		v := src.String()
		n, nErr := strconv.ParseFloat(v, 64)
		if nErr != nil {
			err = fmt.Errorf("copier: float can not support %s source type, src value is not float format string", src.Type().String())
			return
		}
		dst.SetFloat(n)
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := src.Int()
		dst.SetFloat(float64(v))
		break
	case reflect.Float32, reflect.Float64:
		v := src.Float()
		dst.SetFloat(v)
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := src.Uint()
		dst.SetFloat(float64(v))
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
			err = copyFloat(dst, reflect.ValueOf(value))
			return
		}
		err = fmt.Errorf("copier: float not support %s source type", src.Type().String())
		break
	case reflect.Ptr:
		err = copyFloat(dst, src.Elem())
		break
	default:
		err = fmt.Errorf("copier: float not support %s source type", src.Type().String())
		return
	}
	return
}
