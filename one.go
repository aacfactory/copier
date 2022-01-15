package copier

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	tagName = "copy"
)

func copyOne(dst reflect.Value, src reflect.Value) (err error) {
	dst = reflect.Indirect(dst)
	src = reflect.Indirect(src)
	if src.CanConvert(dst.Type()) {
		dst.Set(src.Convert(dst.Type()))
		return
	}

	fieldNum := dst.NumField()
	for i := 0; i < fieldNum; i++ {
		dstFieldValue := dst.Field(i)
		dstFieldType := dst.Type().Field(i)
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

		dstFieldValue0, copyValueErr := copyValue(dstFieldValue, srcFieldValue)
		if copyValueErr != nil {
			err = fmt.Errorf("%s of %s.%s %v", dstFieldType.Name, dst.Type().PkgPath(), dst.Type().Name(), copyValueErr)
			return
		}
		if dstFieldValue0.IsValid() {
			dstFieldValue.Set(dstFieldValue0)
		}
	}
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
	}
	return
}
