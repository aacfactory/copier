package copier

import (
	"database/sql"
	"fmt"
	"reflect"
	"unsafe"
)

var (
	sqlNullStringType  = reflect.TypeOf(sql.NullString{})
	sqlNullInt16Type   = reflect.TypeOf(sql.NullInt16{})
	sqlNullInt32Type   = reflect.TypeOf(sql.NullInt32{})
	sqlNullInt64Type   = reflect.TypeOf(sql.NullInt64{})
	sqlNullFloat64Type = reflect.TypeOf(sql.NullFloat64{})
	sqlNullBoolType    = reflect.TypeOf(sql.NullBool{})
	sqlNullTimeType    = reflect.TypeOf(sql.NullTime{})
)

func copyValue(dst reflect.Value, src reflect.Value) (v reflect.Value, err error) {
	if !dst.CanSet() {
		return
	}

	dstType := dst.Type()
	dstTypeKind := dstType.Kind()
	srcType := src.Type()
	srcTypeKind := srcType.Kind()

	switch dstTypeKind {
	case reflect.String:
		if srcType == sqlNullStringType {
			dst.SetString(src.FieldByName("String").String())
			v = dst
			return
		}
		if srcTypeKind != reflect.String {
			err = fmt.Errorf("type was not matched")
			return
		}
		if !src.CanInterface() {
			return
		}
		dst.SetString(src.String())
		v = dst
		return
	case reflect.Bool:
		if srcType == sqlNullBoolType {
			dst.SetBool(src.FieldByName("Bool").Bool())
			v = dst
			return
		}
		if srcTypeKind != reflect.Bool {
			err = fmt.Errorf("type was not matched")
			return
		}
		if !src.CanInterface() {
			return
		}
		dst.SetBool(src.Bool())
		v = dst
		return
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if srcType == sqlNullInt16Type {
			dst.SetInt(src.FieldByName("Int16").Int())
			v = dst
			return
		}
		if srcType == sqlNullInt32Type {
			dst.SetInt(src.FieldByName("Int32").Int())
			v = dst
			return
		}
		if srcType == sqlNullInt64Type {
			dst.SetInt(src.FieldByName("Int64").Int())
			v = dst
			return
		}
		if srcTypeKind != reflect.Int && srcTypeKind != reflect.Int8 && srcTypeKind != reflect.Int16 && srcTypeKind != reflect.Int32 && srcTypeKind != reflect.Int64 {
			err = fmt.Errorf("type was not matched")
			return
		}
		if !src.CanInterface() {
			return
		}
		dst.SetInt(src.Int())
		v = dst
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
		v = dst
		return
	case reflect.Float32, reflect.Float64:
		if srcType == sqlNullFloat64Type {
			dst.SetFloat(src.FieldByName("Float64").Float())
			v = dst
			return
		}
		if srcTypeKind != reflect.Float32 && srcTypeKind != reflect.Float64 {
			err = fmt.Errorf("type was not matched")
			return
		}
		if !src.CanInterface() {
			return
		}
		dst.SetFloat(src.Float())
		v = dst
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
		v = dst
		return
	case reflect.Struct:
		if srcType == sqlNullTimeType {
			src = src.FieldByName("Time")
			if src.CanConvert(dst.Type()) {
				dst.Set(src.Convert(dst.Type()))
				v = dst
			} else {
				err = fmt.Errorf("type was not matched")
			}
			return
		}
		if srcTypeKind != reflect.Struct {
			err = fmt.Errorf("type was not matched")
			return
		}
		if !src.CanInterface() {
			return
		}
		err = copyStruct(dst, src)
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
		err = copyStruct(reflect.Indirect(dst), reflect.Indirect(src))
		if dst.IsValid() {
			v = dst
		}
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
		dst, err = copyMap(dst, src)
		if dst.IsValid() {
			v = dst
		}
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
