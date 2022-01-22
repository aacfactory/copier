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
		err = fmt.Errorf("map key type is dismatched")
		return
	}
	if src.Len() == 0 {
		return
	}

	srcKeyValues := src.MapKeys()
	for _, srcKeyValue := range srcKeyValues {
		dstKeyValue := srcKeyValue.Convert(dstKeyType)
		srcValValue := src.MapIndex(srcKeyValue)
		if srcValValue.CanConvert(dstValType) {
			dst.SetMapIndex(dstKeyValue, srcValValue.Convert(dstValType))
			continue
		}
		switch srcValValue.Type().Kind() {
		case reflect.Struct:
			if dstValType.Kind() != reflect.Struct {
				err = fmt.Errorf("map value type is dismatched")
				return
			}
			dstValValue := reflect.New(dstValType).Elem()
			vv, cpErr := copyStruct(dstValValue, srcValValue)
			if cpErr != nil {
				err = cpErr
				return
			}
			dst.SetMapIndex(dstKeyValue, vv)
		case reflect.Ptr:
			if dstValType.Kind() != reflect.Ptr {
				err = fmt.Errorf("map value type is dismatched")
				return
			}
			if srcValValue.IsNil() {
				continue
			}

			dstValValue := reflect.New(dstValType.Elem())
			vv, cpErr := copyStruct(dstValValue.Elem(), srcValValue.Elem())
			if cpErr != nil {
				err = cpErr
				return
			}
			dstValValue.Elem().Set(vv)
			dst.SetMapIndex(dstKeyValue, dstValValue)
		case reflect.Array, reflect.Slice:
			if dstValType.Kind() != reflect.Array && dstValType.Kind() != reflect.Slice {
				err = fmt.Errorf("map value type is dismatched")
				return
			}
			if srcValValue.IsNil() || srcValValue.Len() == 0 {
				continue
			}
			dstValValue := reflect.MakeSlice(dstValType.Elem(), 0, 1)
			vv, cpErr := copyArray(dstValValue, srcValValue)
			if cpErr != nil {
				err = cpErr
				return
			}
			dstValValue.Set(vv)
			dst.SetMapIndex(dstKeyValue, dstValValue)
		case reflect.Map:
			if dstValType.Kind() != reflect.Map {
				err = fmt.Errorf("map value type is dismatched")
				return
			}
			dstValValue := reflect.MakeMap(reflect.MapOf(dstValType.Key(), dstValType.Elem()))
			vv, cpErr := copyMap(dstValValue, srcValValue)
			if cpErr != nil {
				err = cpErr
				return
			}
			dst.SetMapIndex(dstKeyValue, vv)
		}
	}
	v = dst
	return
}
