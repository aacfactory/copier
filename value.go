package copier

import (
	"fmt"
	"reflect"
)

func copyValue(dst reflect.Value, src reflect.Value) (err error) {
	srcType := src.Type()
	switch dst.Type().Kind() {
	case reflect.String:
		if srcType.Kind() != reflect.String {
			err = fmt.Errorf("type was not matched")
			return
		}
		dst.Set(src)
	case reflect.Bool:
		if srcType.Kind() != reflect.Bool {
			err = fmt.Errorf("type was not matched")
			return
		}
		dst.Set(src)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		srcTypeKind := srcType.Kind()
		if srcTypeKind == reflect.Int || srcTypeKind == reflect.Int8 || srcTypeKind == reflect.Int16 || srcTypeKind == reflect.Int32 || srcTypeKind == reflect.Int64 {
			dst.SetInt(src.Int())
		} else {
			err = fmt.Errorf("type was not matched")
			return
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		srcTypeKind := srcType.Kind()
		if srcTypeKind == reflect.Uint || srcTypeKind == reflect.Uint8 || srcTypeKind == reflect.Uint16 || srcTypeKind == reflect.Uint32 || srcTypeKind == reflect.Uint64 {
			dst.SetUint(src.Uint())
		} else {
			err = fmt.Errorf("type was not matched")
			return
		}
	case reflect.Float32, reflect.Float64:
		srcTypeKind := srcType.Kind()
		if srcTypeKind == reflect.Float32 || srcTypeKind == reflect.Float64 {
			dst.SetFloat(src.Float())
		} else {
			err = fmt.Errorf("type was not matched")
			return
		}
	case reflect.Complex64, reflect.Complex128:
		srcTypeKind := srcType.Kind()
		if srcTypeKind == reflect.Complex64 || srcTypeKind == reflect.Complex128 {
			dst.SetComplex(src.Complex())
		} else {
			err = fmt.Errorf("type was not matched")
			return
		}
	case reflect.Struct:
		err = copyOne(dst, src)
	case reflect.Ptr:
		if srcType.Kind() == reflect.Ptr {
			if src.IsNil() {
				return
			}
			if dst.Type().AssignableTo()
		}
	case reflect.Array:
		err = copyArray(dst, src)
	case reflect.Map:
		err = copyMap(dst, src)
	default:

	}

	return
}
