package copier

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func copyInt(dst reflect.Value, src reflect.Value) (err error) {
	if src.Type().AssignableTo(dst.Type()) {
		dst.Set(src)
		return
	}
	switch src.Kind() {
	case reflect.String:
		v := src.String()
		n, nErr := strconv.ParseInt(v, 10, 64)
		if nErr != nil {
			err = fmt.Errorf("copier: int can not support %s source type, src value is not float format string", src.Type().String())
			return
		}
		dst.SetInt(n)
		break
	case reflect.Bool:
		v := src.Bool()
		if v {
			dst.SetInt(1)
		} else {
			dst.SetInt(-1)
		}
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := src.Int()
		dst.SetInt(v)
		break
	case reflect.Float32, reflect.Float64:
		v := src.Float()
		dst.SetInt(int64(v))
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := src.Uint()
		dst.SetInt(int64(v))
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
			err = copyInt(dst, reflect.ValueOf(value))
			return
		}
		// time
		if src.Type().ConvertibleTo(timeType) {
			v := src.Interface().(time.Time)
			dst.SetInt(v.UnixMilli())
			return
		}
		err = fmt.Errorf("copier: int not support %s source type", src.Type().String())
		break
	case reflect.Ptr:
		err = copyInt(dst, src.Elem())
		break
	default:
		err = fmt.Errorf("copier: int not support %s source type", src.Type().String())
		return
	}
	return
}
