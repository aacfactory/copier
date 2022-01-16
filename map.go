package copier

import (
	"fmt"
	"reflect"
)

func copyMap(dst reflect.Value, src reflect.Value) (v reflect.Value, err error) {
	dstType := dst.Type()
	dstKeyType := dstType.Key()
	dstValType := dstType.Elem()

	srcType := src.Type()
	srcKeyType := srcType.Key()

	if !srcKeyType.ConvertibleTo(dstKeyType) {
		err = fmt.Errorf("key type was not matched")
		return
	}
	if src.Len() == 0 {
		return
	}

	dst = reflect.MakeMap(reflect.MapOf(dstKeyType, dstValType))

	srcKeyValues := src.MapKeys()
	for _, srcKeyValue := range srcKeyValues {
		dstKeyValue := srcKeyValue.Convert(dstKeyType)
		srcValValue := src.MapIndex(srcKeyValue)
		if srcValValue.CanConvert(dstValType) {
			dst.SetMapIndex(dstKeyValue, srcValValue.Convert(dstValType))
			continue
		}
		var dstValValue reflect.Value
		switch srcValValue.Type().Kind() {
		case reflect.Struct:
			if dstValType.Kind() != reflect.Struct {
				err = fmt.Errorf("key type was not matched")
				return
			}
			dstValValue = reflect.Indirect(reflect.New(dstValType))
			dstValValue, err = copyValue(dstValValue, srcValValue)
			if err != nil {
				return
			}
		case reflect.Ptr:
			if dstValType.Kind() != reflect.Ptr {
				err = fmt.Errorf("key type was not matched")
				return
			}
			dstValValue = reflect.New(dstValType)
			dstValValue, err = copyValue(dstValValue, srcValValue)
			if err != nil {
				return
			}
		case reflect.Array, reflect.Slice:
			if dstValType.Kind() != reflect.Array && dstValType.Kind() != reflect.Slice {
				err = fmt.Errorf("key type was not matched")
				return
			}
			dstValValue = reflect.MakeSlice(dstValType.Elem(), 0, 1)
			dstValValue, err = copyArray(dstValValue, srcValValue)
			if err != nil {
				return
			}
		case reflect.Map:
			if dstValType.Kind() != reflect.Map {
				err = fmt.Errorf("key type was not matched")
				return
			}
			dstValValue = reflect.MakeMap(reflect.MapOf(dstValType.Key(), dstValType.Elem()))
			dstValValue, err = copyMap(dstValValue, srcValValue)
			if err != nil {
				return
			}
		}
		dst.SetMapIndex(dstKeyValue, dstValValue)
	}
	v = dst
	return
}
