package copier

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"time"
)

func copyTime(dst reflect.Value, src reflect.Value) (err error) {
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
		t, parseErr := time.Parse(time.RFC3339, v)
		if parseErr != nil {
			err = parseErr
			return
		}
		dst.Set(reflect.ValueOf(t))
		break
	case reflect.Int, reflect.Int32, reflect.Int64:
		v := src.Int()
		if v < 1 {
			return
		}
		t := time.UnixMilli(v)
		dst.Set(reflect.ValueOf(t))
		break
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := src.Uint()
		if v < 1 {
			return
		}
		t := time.UnixMilli(int64(v))
		dst.Set(reflect.ValueOf(t))
		break
	case reflect.Struct:
		// time
		if src.Type().ConvertibleTo(timeType) {
			v := src.Interface().(time.Time)
			dst.SetString(v.Format(time.RFC3339))
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
			err = copyTime(dst, reflect.ValueOf(value))
			return
		}
		err = fmt.Errorf("copier: time not support %s source type", src.Type().String())
		break
	case reflect.Ptr:
		err = copyTime(dst, src.Elem())
		break
	default:
		err = fmt.Errorf("copier: time not support %s source type", src.Type().String())
		return
	}
	return
}
