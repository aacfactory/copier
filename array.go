package copier

import (
	"fmt"
	"reflect"
)

func copyArray(dst reflect.Value, src reflect.Value) (v reflect.Value, err error) {
	if src.IsNil() || src.Len() == 0 {
		return
	}

	dstElemType := dst.Type().Elem()
	dstElemTypeKind := dstElemType.Kind()
	srcElemType := src.Type().Elem()
	srcElemTypeKind := srcElemType.Kind()

	// bytes
	if dstElemTypeKind == reflect.Uint8 {
		if srcElemTypeKind == reflect.Uint8 {
			dst = src.Convert(dst.Type())
			v = dst
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
			vv, cpErr := copyStruct(dstItem.Elem(), srcItem.Elem())
			if cpErr != nil {
				err = cpErr
				return
			}
			dstItem.Elem().Set(vv)
			dst = reflect.Append(dst, dstItem)
		} else if dstElemTypeKind == reflect.Struct {
			dstItem := reflect.New(dstElemType).Elem()
			vv, cpErr := copyStruct(dstItem, srcItem)
			if cpErr != nil {
				err = cpErr
				return
			}
			dstItem.Set(vv)
			dst = reflect.Append(dst, dstItem)
		} else if dstElemTypeKind == reflect.Slice {
			dstItem := reflect.MakeSlice(dstElemType, 0, 1)
			vv, cpErr := copyArray(dstItem, srcItem)
			if cpErr != nil {
				err = cpErr
				return
			}
			if !dstItem.IsValid() {
				continue
			}
			dst = reflect.Append(dst, vv)
		}

	}
	v = dst
	return
}
