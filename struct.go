package copier

import (
	"reflect"
)

const (
	tagKey     = "copy"
	discardTag = "-"
)

func copyStruct(dst reflect.Value, src reflect.Value) (err error) {
	if src.Type().AssignableTo(dst.Type()) {
		dst.Set(src)
		return
	}
	if src.Kind() == reflect.Ptr {
		src = src.Elem()
	}
	dstType := dst.Type()
	// time
	if dstType.ConvertibleTo(timeType) {
		err = copyTime(dst, src)
		return
	}
	// sql
	if dstType.Implements(sqlValuerType) {
		err = copySQLScanner(dst, src)
		return
	}
	// addr
	if dst.CanAddr() {
		dstAddr := dst.Addr()
		// sql
		if dstType.Implements(sqlScannerType) {
			err = copySQLScanner(dst, src)
			return
		}
		// text
		if dstAddr.Type().Implements(textUnmarshalerType) {
			err = copyValue(dstAddr, src)
			return
		}
	}
	srcType := src.Type()
	// fields
	srcFieldNum := srcType.NumField()
	if srcFieldNum == 0 {
		return
	}
	fieldNum := dstType.NumField()
	if fieldNum == 0 {
		return
	}
	for i := 0; i < fieldNum; i++ {
		dstFieldType := dstType.Field(i)
		if dstFieldType.Anonymous {
			if dstFieldType.Type.Kind() == reflect.Ptr {
				if !dstFieldType.IsExported() {
					continue
				}
				dstField := reflect.New(dstFieldType.Type.Elem())
				if err = copyStruct(dstField.Elem(), src); err != nil {
					return
				}
				dst.FieldByName(dstFieldType.Name).Set(dstField)
				continue
			}
			dstField := dst.FieldByName(dstFieldType.Name)
			if err = copyStruct(dstField, src); err != nil {
				return
			}
			continue
		}
		if !dstFieldType.IsExported() {
			continue
		}
		key := ""
		tag, hasTag := dstFieldType.Tag.Lookup(tagKey)
		if hasTag {
			if tag == discardTag {
				continue
			}
			for j := 0; j < srcFieldNum; j++ {
				srcField := srcType.Field(j)
				srcTag, hasSrcTag := srcField.Tag.Lookup(tagKey)
				if hasSrcTag && tag == srcTag {
					key = srcField.Name
					break
				}
			}
		} else {
			key = dstFieldType.Name
		}
		if key == "" {
			continue
		}
		srcFieldType, hasSrcField := srcType.FieldByName(key)
		if !hasSrcField {
			continue
		}
		srcField := src.FieldByName(srcFieldType.Name)
		dstField := dst.FieldByName(dstFieldType.Name)
		if err = copyValue(dstField, srcField); err != nil {
			return
		}
	}
	return
}
