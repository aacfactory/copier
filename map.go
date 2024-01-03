package copier

import (
	"reflect"
)

func copyMap(dst reflect.Value, src reflect.Value) (err error) {
	if src.Type().AssignableTo(dst.Type()) {
		dst.Set(src)
		return
	}
	dstType := dst.Type()
	dstKeyType := dstType.Key()
	dstValType := dstType.Elem()
	srcIter := src.MapRange()
	for srcIter.Next() {
		// key
		srcKey := srcIter.Key()
		dstKey := reflect.New(dstKeyType).Elem()
		if err = copyValue(dstKey, srcKey); err != nil {
			return
		}
		// val
		srcVal := srcIter.Value()
		dstVal := reflect.New(dstValType).Elem()
		if err = copyValue(dstVal, srcVal); err != nil {
			return
		}
		dst.SetMapIndex(dstKey, dstVal)
	}
	return
}
