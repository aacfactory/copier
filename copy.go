package copier

import (
	"errors"
	"fmt"
	"reflect"
)

func Copy(dst any, src any) (err error) {
	if src == nil {
		return
	}
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr {
		err = fmt.Errorf("copier: dst must be ptr")
		return
	}
	srcValue := reflect.ValueOf(src)
	err = copyValue(dstValue.Elem(), srcValue)
	if err != nil {
		err = errors.Join(fmt.Errorf("copier: copy failed"), err)
		return
	}
	return
}

func ValueOf[E any](src any) (dst E, err error) {
	p := new(E)
	if err = Copy(p, src); err != nil {
		return
	}
	dst = *p
	return
}
