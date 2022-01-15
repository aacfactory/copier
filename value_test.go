package copier_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

func TestTime(_ *testing.T) {

	x := 10

	fmt.Println(reflect.ValueOf(x).CanConvert(reflect.TypeOf(float64(10.10))))

	tt := time.Now()
	fmt.Println(reflect.ValueOf(tt).CanConvert(reflect.TypeOf(time.Time{})))

	fmt.Println(reflect.TypeOf(Foo{}).Field(3).Type.Kind())

	t := &Foo{}
	//d := Date(time.Now())
	d := &Foo{}

	tv := reflect.ValueOf(t).Elem()
	dv := reflect.ValueOf(d).Elem()

	fmt.Println(tv.CanAddr(), tv.CanSet(), tv.CanInterface(), dv.CanInterface())
	tv.FieldByName("Str").Set(reflect.ValueOf("x"))
	tv.FieldByName("Time").Set(reflect.ValueOf(time.Now()).Convert(reflect.TypeOf(Date{})))
	SetUnexportedValue(tv.FieldByName("i"), complex128(10))
	fmt.Println(tv.Interface())

	fmt.Println(reflect.TypeOf(Date{}).NumField())
	fmt.Println(tv.FieldByName("i").CanAddr(), tv.FieldByName("i").CanSet(), tv.FieldByName("i").CanInterface())
}

func GetUnexportedValue(value reflect.Value) interface{} {
	return reflect.NewAt(value.Type(), unsafe.Pointer(value.UnsafeAddr())).Elem().Interface()
}

func SetUnexportedValue(value reflect.Value, v interface{}) {
	rv := reflect.ValueOf(v)
	if rv.CanConvert(value.Type()) {
		rv = rv.Convert(value.Type())
		reflect.NewAt(value.Type(), unsafe.Pointer(value.UnsafeAddr())).Elem().Set(rv)
	}
}
