package copier

import (
	"fmt"
	"reflect"
)

func copyArray(dst reflect.Value, src reflect.Value) (v reflect.Value, err error) {
	if !dst.CanSet() {
		return
	}
	if src.IsNil() {
		return
	}

	dstElemType := dst.Type().Elem()
	dstElemTypeKind := dstElemType.Kind()
	srcElemType := src.Type().Elem()
	srcElemTypeKind := srcElemType.Kind()

	// bytes
	if dstElemTypeKind == reflect.Uint8 {
		if srcElemTypeKind == reflect.Uint8 {
			dst.SetBytes(src.Bytes())
			return
		} else {
			err = fmt.Errorf("type was not matched")
			return
		}
	}
	// sample
	if srcElemType.AssignableTo(dstElemType) {
		v = reflect.AppendSlice(dst, src)
		return
	}
	// copy elem
	size := src.Len()
	for i := 0; i < size; i++ {
		srcItem := src.Index(i)
		if dstElemTypeKind == reflect.Ptr {
			dstItem := reflect.New(dstElemType.Elem())
			dstItem, err = copyValue(dstItem, srcItem)
			if err != nil {
				return
			}
			if !dstItem.IsValid() {
				continue
			}
			dst = reflect.Append(dst, dstItem)
		} else if dstElemTypeKind == reflect.Struct {
			dstItem := reflect.New(dstElemType)
			err = copyOne(dstItem, srcItem)
			if err != nil {
				return
			}
			if !dstItem.IsValid() {
				continue
			}
			dst = reflect.Append(dst, dstItem.Elem())
		} else if dstElemTypeKind == reflect.Slice {
			dstItem := reflect.MakeSlice(dstElemType, 0, 1)
			dstItem, err = copyArray(dstItem, srcItem)
			if err != nil {
				return
			}
			if !dstItem.IsValid() {
				continue
			}
			dst = reflect.Append(dst, dstItem)
		}

	}
	v = dst
	return
}
