package copier

import (
	"reflect"
	"strings"
)

const (
	tagKey     = "copy"
	discardTag = '-'
	comma      = ','
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
		srcFieldIdx := -1
		setFn := ""
		getFn := ""
		tag, hasTag := dstFieldType.Tag.Lookup(tagKey)
		if hasTag {
			tag = strings.TrimSpace(tag)
			if strings.IndexByte(tag, discardTag) == 0 {
				continue
			}
			// set func of dst
			if idx := strings.IndexByte(tag, comma); idx > -1 {
				setFn = strings.TrimSpace(tag[idx+1:])
				tag = strings.TrimSpace(tag[0:idx])
			}
			for j := 0; j < srcFieldNum; j++ {
				srcField := srcType.Field(j)
				srcTag, hasSrcTag := srcField.Tag.Lookup(tagKey)
				if hasSrcTag {
					srcTag = strings.TrimSpace(srcTag)
					if strings.IndexByte(srcTag, discardTag) == 0 {
						continue
					}
					if idx := strings.IndexByte(srcTag, comma); idx > -1 {
						getFn = strings.TrimSpace(srcTag[idx+1:])
						srcTag = strings.TrimSpace(srcTag[0:idx])
						if srcTag == tag || srcTag == dstFieldType.Name || srcField.Name == dstFieldType.Name || srcField.Name == tag {
							srcFieldIdx = j
							break
						}
						getFn = ""
					}
					continue
				}
				if srcFieldIdx > -1 || getFn != "" {
					break
				}
				if srcField.Name == tag {
					srcFieldIdx = j
					break
				}
			}
			if srcFieldIdx == -1 && getFn == "" {
				setFn = ""
			}
		} else {
			for j := 0; j < srcFieldNum; j++ {
				srcField := srcType.Field(j)
				if srcField.Name == dstFieldType.Name {
					srcFieldIdx = j
					break
				}
			}
		}
		var srcField reflect.Value
		if srcFieldIdx == -1 && getFn == "" {
			continue
		}
		if getFn != "" {
			var getMethod reflect.Value
			_, hasGetMethod := srcType.MethodByName(getFn)
			if hasGetMethod {
				getMethod = src.MethodByName(getFn)
			} else {
				continue
			}
			srcField = getMethod.Call(nil)[0]
		} else if srcFieldIdx > -1 {
			srcField = src.Field(srcFieldIdx)
		} else {
			continue
		}
		if setFn != "" && dstFieldType.Type != srcField.Type() {
			dst.Addr().MethodByName(setFn).Call([]reflect.Value{srcField})
		} else {
			dstField := dst.Field(i)
			if err = copyValue(dstField, srcField); err != nil {
				return
			}
		}
	}
	return
}
