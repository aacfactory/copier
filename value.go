package copier

import (
	"fmt"
	"reflect"
	"unsafe"
)

func copyValue(dst reflect.Value, src reflect.Value) (v reflect.Value, err error) {
	if !dst.CanSet() {
		return
	}

	dstType := src.Type()
	dstTypeKind := dstType.Kind()
	srcType := src.Type()
	srcTypeKind := srcType.Kind()

	switch dstTypeKind {
	case reflect.String:
		if srcTypeKind != reflect.String {
			err = fmt.Errorf("type was not matched")
			return
		}
		if !src.CanInterface() {
			return
		}
		dst.SetString(src.String())
		return
	case reflect.Bool:
		if srcTypeKind != reflect.Bool {
			err = fmt.Errorf("type was not matched")
			return
		}
		if !src.CanInterface() {
			return
		}
		dst.SetBool(src.Bool())
		return
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if srcTypeKind != reflect.Int && srcTypeKind != reflect.Int8 && srcTypeKind != reflect.Int16 && srcTypeKind != reflect.Int32 && srcTypeKind != reflect.Int64 {
			err = fmt.Errorf("type was not matched")
			return
		}
		if !src.CanInterface() {
			return
		}
		dst.SetInt(src.Int())
		return
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if srcTypeKind != reflect.Uint && srcTypeKind != reflect.Uint8 && srcTypeKind != reflect.Uint16 && srcTypeKind != reflect.Uint32 && srcTypeKind != reflect.Uint64 {
			err = fmt.Errorf("type was not matched")
			return
		}
		if !src.CanInterface() {
			return
		}
		dst.SetUint(src.Uint())
		return
	case reflect.Float32, reflect.Float64:
		if srcTypeKind != reflect.Float32 && srcTypeKind != reflect.Float64 {
			err = fmt.Errorf("type was not matched")
			return
		}
		if !src.CanInterface() {
			return
		}
		dst.SetFloat(src.Float())
		return
	case reflect.Complex64, reflect.Complex128:
		if srcTypeKind != reflect.Complex64 && srcTypeKind != reflect.Complex128 {
			err = fmt.Errorf("type was not matched")
			return
		}
		if !src.CanInterface() {
			return
		}
		dst.SetComplex(src.Complex())
		return
	case reflect.Struct:
		if srcTypeKind != reflect.Struct {
			err = fmt.Errorf("type was not matched")
			return
		}
		if !src.CanInterface() {
			return
		}
		err = copyOne(dst, src)
		v = dst
		return
	case reflect.Ptr:
		if srcTypeKind != reflect.Ptr {
			err = fmt.Errorf("type was not matched")
			return
		}
		if src.IsNil() {
			return
		}
		if !src.CanAddr() {
			return
		}
		err = copyOne(reflect.Indirect(dst), reflect.Indirect(src))
		v = dst
		return
	case reflect.Array, reflect.Slice:
		if srcTypeKind != reflect.Array && srcTypeKind != reflect.Slice {
			err = fmt.Errorf("type was not matched")
			return
		}
		dst, err = copyArray(dst, src)
		if dst.IsValid() {
			v = dst
		}
		return
	case reflect.Map:
		if srcTypeKind != reflect.Map {
			err = fmt.Errorf("type was not matched")
			return
		}
		if src.IsNil() {
			return
		}
		err = copyMap(dst, src)
		return
	}
	return
}

func getUnexportedValue(value reflect.Value) interface{} {
	return reflect.NewAt(value.Type(), unsafe.Pointer(value.UnsafeAddr())).Elem().Interface()
}

func setUnexportedValue(value reflect.Value, v interface{}) {
	rv := reflect.ValueOf(v)
	if rv.CanConvert(value.Type()) {
		rv = rv.Convert(value.Type())
		reflect.NewAt(value.Type(), unsafe.Pointer(value.UnsafeAddr())).Elem().Set(rv)
	}
}

func isNumber(v reflect.Kind) (ok bool) {
	isInt := v == reflect.Int || v == reflect.Int8 || v == reflect.Int16 || v == reflect.Int32 || v == reflect.Int64
	isUint := v == reflect.Uint || v == reflect.Uint8 || v == reflect.Uint16 || v == reflect.Uint32 || v == reflect.Uint64
	isFloat := v == reflect.Float32 || v == reflect.Float64
	ok = isInt || isUint || isFloat
	return
}
