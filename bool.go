package copier

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
)

const (
	trueString       = "true"
	trueUpperString  = "TRUE"
	falseString      = "false"
	falseUpperString = "FALSE"
)

func copyBool(dst reflect.Value, src reflect.Value) (err error) {
	if src.Type().AssignableTo(dst.Type()) {
		dst.Set(src)
		return
	}
	switch src.Kind() {
	case reflect.String:
		v := src.String()
		if v == trueString || strings.ToUpper(v) == trueUpperString {
			dst.SetBool(true)
		} else if v == falseString || strings.ToUpper(v) == falseUpperString {
			dst.SetBool(false)
		} else {
			err = fmt.Errorf("copier: %s source type can not convert to bool", src.Type().String())
			return
		}
		break
	case reflect.Bool:
		dst.Set(src)
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := src.Int()
		if v == 1 {
			dst.SetBool(true)
		} else if v == 0 || v == -1 {
			dst.SetBool(false)
		} else {
			err = fmt.Errorf("copier: %s source type can not convert to bool", src.Type().String())
			return
		}
		break
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := src.Uint()
		if v == 1 {
			dst.SetBool(true)
		} else if v == 0 {
			dst.SetBool(false)
		} else {
			err = fmt.Errorf("copier: %s source type can not convert to bool", src.Type().String())
			return
		}
		break
	case reflect.Uint8:
		v := src.Uint()
		if v == 'T' || v == 't' {
			dst.SetBool(true)
		} else if v == 'F' || v == 'f' {
			dst.SetBool(false)
		} else {
			err = fmt.Errorf("copier: %s source type can not convert to bool", src.Type().String())
			return
		}
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
			err = copyBool(dst, reflect.ValueOf(value))
			return
		}
		err = fmt.Errorf("copier: bool not support %s source type", src.Type().String())
		break
	case reflect.Ptr:
		err = copyBool(dst, src.Elem())
		break
	default:
		err = fmt.Errorf("copier: bool not support %s source type", src.Type().String())
		return
	}
	return
}
