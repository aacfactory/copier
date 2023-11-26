package copier

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	tagName = "copy"
)

func copyStruct(dst reflect.Value, src reflect.Value) (v reflect.Value, err error) {
	dstType := dst.Type()
	srcType := src.Type()
	if srcType.Kind() != reflect.Struct {
		err = fmt.Errorf("type is dismatched")
		return
	}

	if srcType == dstType {
		dst.Set(src)
		v = dst
		return
	}

	if src.CanConvert(dstType) {
		dst.Set(src.Convert(dstType))
		v = dst
		return
	}

	if srcType.Implements(sqlScannerType) || reflect.New(srcType).Type().Implements(sqlScannerType) {
		var srcValue reflect.Value
		valid := false
		srcFields := src.NumField()
		for i := 0; i < srcFields; i++ {
			srcFieldType := srcType.Field(i)
			if srcFieldType.Anonymous {
				continue
			}
			if srcFieldType.Name == "Valid" {
				valid = src.Field(i).Bool()
				continue
			}
			srcValue = src.Field(i)
		}
		if valid {
			if srcValue.Type() == dstType {
				dst.Set(srcValue)
				v = dst
				return
			}
			if src.CanConvert(dstType) {
				dst.Set(src.Convert(dstType))
				v = dst
				return
			}
		}
		return
	}

	fieldNum := dst.NumField()
	for i := 0; i < fieldNum; i++ {
		dstFieldValue := dst.Field(i)
		dstFieldType := dstType.Field(i)
		// Anonymous
		if dstFieldType.Anonymous {
			vv, cpErr := copyStruct(dstFieldValue, src)
			if cpErr != nil {
				err = cpErr
				return
			}
			dstFieldValue.Set(vv)
			continue
		}
		tag, hasTag := dstFieldType.Tag.Lookup(tagName)
		var srcFieldValue reflect.Value
		found := false
		if hasTag {
			tag = strings.TrimSpace(tag)
			if tag == "-" {
				continue
			}
			srcFieldValue, found = findFieldValueByTag(tag, src)
		} else {
			srcFieldValue, found = findFieldValueByName(dstFieldType.Name, src)
		}
		if !found {
			continue
		}

		srcFieldType := srcFieldValue.Type()
		srcFieldTypeKind := srcFieldType.Kind()
		switch dstFieldType.Type.Kind() {
		case reflect.String:
			if srcFieldType == sqlNullStringType {
				if srcFieldValue.FieldByName("Valid").Bool() {
					dstFieldValue.SetString(srcFieldValue.FieldByName("String").String())
				}
				continue
			}
			if srcFieldTypeKind != reflect.String {
				err = fmt.Errorf("type is dismatched")
				return
			}
			if srcFieldValue.CanInterface() {
				dstFieldValue.SetString(srcFieldValue.String())
			}
		case reflect.Bool:
			if srcFieldType == sqlNullBoolType {
				if srcFieldValue.FieldByName("Valid").Bool() {
					dstFieldValue.SetBool(src.FieldByName("Bool").Bool())
				}
				continue
			}
			if srcFieldTypeKind != reflect.Bool {
				err = fmt.Errorf("type is dismatched")
				return
			}
			if srcFieldValue.CanInterface() {
				dstFieldValue.SetBool(srcFieldValue.Bool())
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if srcFieldType == sqlNullInt16Type {
				if srcFieldValue.FieldByName("Valid").Bool() {
					dstFieldValue.SetInt(srcFieldValue.FieldByName("Int16").Int())
				}
				continue
			}
			if srcFieldType == sqlNullInt32Type {
				if srcFieldValue.FieldByName("Valid").Bool() {
					dstFieldValue.SetInt(srcFieldValue.FieldByName("Int32").Int())
				}
				continue
			}
			if srcFieldType == sqlNullInt64Type {
				if srcFieldValue.FieldByName("Valid").Bool() {
					dstFieldValue.SetInt(srcFieldValue.FieldByName("Int64").Int())
				}
				continue
			}
			if srcFieldTypeKind != reflect.Int && srcFieldTypeKind != reflect.Int8 && srcFieldTypeKind != reflect.Int16 && srcFieldTypeKind != reflect.Int32 && srcFieldTypeKind != reflect.Int64 {
				err = fmt.Errorf("type is dismatched")
				return
			}
			if srcFieldValue.CanInterface() {
				dstFieldValue.SetInt(srcFieldValue.Int())
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if srcFieldTypeKind != reflect.Uint && srcFieldTypeKind != reflect.Uint8 && srcFieldTypeKind != reflect.Uint16 && srcFieldTypeKind != reflect.Uint32 && srcFieldTypeKind != reflect.Uint64 {
				err = fmt.Errorf("type is dismatched")
				return
			}
			if srcFieldValue.CanInterface() {
				dstFieldValue.SetUint(srcFieldValue.Uint())
			}
		case reflect.Float32, reflect.Float64:
			if srcFieldType == sqlNullFloat64Type {
				if srcFieldValue.FieldByName("Valid").Bool() {
					dstFieldValue.SetFloat(srcFieldValue.FieldByName("Float64").Float())
				}
				continue
			}
			if srcFieldTypeKind != reflect.Float32 && srcFieldTypeKind != reflect.Float64 {
				err = fmt.Errorf("type is dismatched")
				return
			}
			if srcFieldValue.CanInterface() {
				dstFieldValue.SetFloat(srcFieldValue.Float())
			}
		case reflect.Complex64, reflect.Complex128:
			if srcFieldTypeKind != reflect.Complex64 && srcFieldTypeKind != reflect.Complex128 {
				err = fmt.Errorf("type is dismatched")
				return
			}
			if srcFieldValue.CanInterface() {
				dstFieldValue.SetComplex(srcFieldValue.Complex())
			}
		case reflect.Struct:
			if srcFieldTypeKind != reflect.Struct {
				err = fmt.Errorf("type is dismatched")
				return
			}
			vv, cpErr := copyStruct(dstFieldValue, srcFieldValue)
			if cpErr != nil {
				err = cpErr
				return
			}
			dstFieldValue.Set(vv)
		case reflect.Ptr:
			if srcFieldTypeKind != reflect.Ptr {
				err = fmt.Errorf("type is dismatched")
				return
			}
			if srcFieldValue.IsNil() {
				continue
			}
			dstFieldValueValue := reflect.New(dstFieldType.Type.Elem())
			vv, cpErr := copyStruct(dstFieldValueValue.Elem(), srcFieldValue.Elem())
			if cpErr != nil {
				err = cpErr
				return
			}
			dstFieldValueValue.Elem().Set(vv)
			dstFieldValue.Set(dstFieldValueValue)
		case reflect.Array, reflect.Slice:
			if srcFieldTypeKind != reflect.Array && srcFieldTypeKind != reflect.Slice {
				err = fmt.Errorf("type is dismatched")
				return
			}
			if srcFieldValue.IsNil() || srcFieldValue.Len() == 0 {
				continue
			}
			dstFieldValueValue := reflect.MakeSlice(dstFieldType.Type, 0, 1)
			vv, cpErr := copyArray(dstFieldValueValue, srcFieldValue)
			if cpErr != nil {
				err = cpErr
				return
			}
			dstFieldValue.Set(vv)
		case reflect.Map:
			if srcFieldTypeKind != reflect.Map {
				err = fmt.Errorf("type is dismatched")
				return
			}
			if srcFieldValue.IsNil() || srcFieldValue.Len() == 0 {
				continue
			}
			dstFieldValueValue := reflect.MakeMap(reflect.MapOf(dstFieldType.Type.Key(), dstFieldType.Type.Elem()))
			vv, cpErr := copyMap(dstFieldValueValue, srcFieldValue)
			if cpErr != nil {
				err = cpErr
				return
			}
			dstFieldValue.Set(vv)
		}
	}

	v = dst
	return
}

func findFieldValueByTag(tag string, src reflect.Value) (v reflect.Value, has bool) {
	srcType := src.Type()
	fieldNum := srcType.NumField()
	for i := 0; i < fieldNum; i++ {
		srcFieldType := srcType.Field(i)
		srcTag, hasTag := srcFieldType.Tag.Lookup(tagName)
		if hasTag && strings.TrimSpace(srcTag) == tag {
			v = src.Field(i)
			has = true
			return
		}
	}
	return
}

func findFieldValueByName(name string, src reflect.Value) (v reflect.Value, has bool) {
	if _, has = src.Type().FieldByName(name); has {
		v = src.FieldByName(name)
		return
	}
	srcType := src.Type()
	fields := srcType.NumField()
	for i := 0; i < fields; i++ {
		srcFieldType := srcType.Field(i)
		if srcFieldType.Anonymous {
			v, has = findFieldValueByName(name, src.Field(i))
			if has {
				return
			}
		}
	}
	return
}
