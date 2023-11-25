package copier

import (
	"fmt"
	"reflect"
)

func Copy(dst any, src any) (err error) {
	if dst == nil {
		err = fmt.Errorf("copy failed for dst is nil")
		return
	}
	if src == nil {
		return
	}
	dstType := reflect.TypeOf(dst)
	if dstType.Kind() != reflect.Ptr {
		err = fmt.Errorf("copy failed for type of dst must be ptr")
		return
	}
	dstType = dstType.Elem()
	dstValue := reflect.ValueOf(dst)
	srcValue := reflect.ValueOf(src)
	srcType := reflect.TypeOf(src)
	if srcType.Kind() == reflect.Ptr {
		srcType = srcType.Elem()
		srcValue = srcValue.Elem()
	}
	if dstType.Kind() != srcType.Kind() {
		err = fmt.Errorf("copy failed for type between dst and src is not matched")
		return
	}
	switch dstType.Kind() {
	case reflect.Struct:
		cpValue, cpErr := copyStruct(dstValue.Elem(), srcValue)
		if cpErr != nil {
			err = fmt.Errorf("copy failed for %v", cpErr)
			return
		}
		dstValue.Elem().Set(cpValue)
	case reflect.Array, reflect.Slice:
		cpValue, cpErr := copyArray(dstValue.Elem(), srcValue)
		if cpErr != nil {
			err = fmt.Errorf("copy failed for %v", cpErr)
			return
		}
		dstValue.Elem().Set(cpValue)
	case reflect.Map:
		cpValue, cpErr := copyMap(dstValue.Elem(), srcValue)
		if cpErr != nil {
			err = fmt.Errorf("copy failed for %v", cpErr)
			return
		}
		dstValue.Elem().Set(cpValue)
	default:
		err = fmt.Errorf("copy failed for %v is not supported", dstType.Kind())
		return
	}
	if err != nil {
		err = fmt.Errorf("copy failed for %v", err)
		return
	}
	return
}

func ValueOf[D any](src any) (dst D, err error) {
	err = Copy(&dst, src)
	return
}
