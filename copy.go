package copier

import (
	"fmt"
	"github.com/modern-go/reflect2"
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

	dstType := reflect2.TypeOf(dst)
	if dstType.Kind() != reflect.Ptr {
		err = fmt.Errorf("copy failed for type of dst must be ptr")
		return
	}
	// dst element type
	dstType = dstType.(reflect2.PtrType).Elem()
	// dst ptr
	dstPtr := reflect2.PtrOf(dst)
	// src type
	srcType := reflect2.TypeOf(src)
	if srcType.Kind() == reflect.Ptr {
		srcType = srcType.(reflect2.PtrType).Elem()
	}
	// src ptr
	srcPtr := reflect2.PtrOf(src)
	// copy by type
	switch dstType.Kind() {
	case reflect.String:

		break
	case reflect.Bool:

		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		break
	case reflect.Float32, reflect.Float64:

		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

		break
	case reflect.Struct:
		break
	case reflect.Slice:
		break
	case reflect.Map:
		break
	default:
		err = fmt.Errorf("copy failed for %v is not supported", dstType.Kind())
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
