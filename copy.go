package copier

import (
	"fmt"
	"reflect"
)

func Copy(dst interface{}, src interface{}) (err error) {
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
	dstValue := reflect.Indirect(reflect.ValueOf(dst))
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
		err = copyOne(dstValue, srcValue)
	case reflect.Array, reflect.Slice:
		err = copyArray(dstValue, srcValue)
	case reflect.Map:
		err = copyMap(dstValue, srcValue)
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
