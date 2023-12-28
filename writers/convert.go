package writers

import (
	"github.com/modern-go/reflect2"
	"reflect"
	"unsafe"
)

type Convertible interface {
	Convert() any
}

func IsConvertible(typ reflect2.Type) bool {
	return typ.Implements(convertibleType)
}

func convert(ptr unsafe.Pointer, typ reflect2.Type) (unsafe.Pointer, reflect2.Type) {
	st := typ.Type1()
	if st.Kind() == reflect.Ptr {
		st = st.Elem()
	}
	convertible := reflect.NewAt(st, ptr).Interface().(Convertible)
	v := convertible.Convert()
	return reflect2.PtrOf(v), reflect2.TypeOf(v)
}
