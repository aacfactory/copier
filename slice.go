package copier

import (
	"fmt"
	"reflect"
)

func copySlice(dst reflect.Value, src reflect.Value) (err error) {
	if dst.Type().Elem().Kind() == reflect.Uint8 {
		err = copyBytes(dst, src)
		return
	}
	if src.Kind() == reflect.Ptr {
		src = src.Elem()
	}
	if src.Kind() != reflect.Slice {
		err = fmt.Errorf("copier: slice can not support %s source type", src.Type().String())
		return
	}
	srcLen := src.Len()
	if srcLen == 0 {
		return
	}
	dstValue := reflect.MakeSlice(dst.Type(), 0, srcLen)
	dstElemType := dstValue.Type().Elem()
	for i := 0; i < srcLen; i++ {
		srcElem := src.Index(i)
		dstElem := reflect.New(dstElemType).Elem()
		if err = copyValue(dstElem, srcElem); err != nil {
			return
		}
		dstValue = reflect.Append(dstValue, dstElem)
	}
	dst.Set(dstValue)
	return
}
