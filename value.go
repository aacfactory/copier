package copier

import (
	"fmt"
	"reflect"
)

func copyValue(dst reflect.Value, src reflect.Value) (err error) {
	if src.Type().AssignableTo(dst.Type()) {
		dst.Set(src)
		return
	}
	switch dst.Kind() {
	case reflect.String:
		err = copyString(dst, src)
		break
	case reflect.Bool:
		err = copyBool(dst, src)
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		err = copyInt(dst, src)
		break
	case reflect.Float32, reflect.Float64:
		err = copyFloat(dst, src)
		break
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		err = copyUint(dst, src)
		break
	case reflect.Uint8:
		err = copyByte(dst, src)
		break
	case reflect.Struct:
		err = copyStruct(dst, src)
		break
	case reflect.Ptr:
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		// text
		if dst.Type().Implements(textUnmarshalerType) {
			err = copyText(dst, src)
			return
		}
		err = copyValue(dst.Elem(), src)
		break
	case reflect.Slice:
		if dst.IsNil() {
			dst.Set(reflect.MakeSlice(dst.Type(), 0, 1))
		}
		err = copySlice(dst, src)
		break
	case reflect.Map:
		if dst.IsNil() {
			dst.Set(reflect.MakeMap(dst.Type()))
		}
		err = copyMap(dst, src)
		break
	default:
		err = fmt.Errorf("copier: %s is not support", dst.Type().String())
		return
	}
	return
}
