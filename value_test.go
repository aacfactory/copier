package copier_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

func TestTime(_ *testing.T) {

	t := &Foo{}
	//d := Date(time.Now())
	d := &Foo{}

	tv := reflect.ValueOf(t).Elem()
	dv := reflect.ValueOf(d).Elem()

	fmt.Println(tv.CanAddr(), tv.CanSet(), tv.CanInterface(), dv.CanInterface())
	tv.FieldByName("Str").Set(reflect.ValueOf("x"))
	tv.FieldByName("Time").Set(reflect.ValueOf(time.Now()).Convert(reflect.TypeOf(Date{})))
	tv.FieldByName("i").SetInt(int64(10))
	tv.FieldByName("i").Type().
		fmt.Println(t)
	fmt.Println(tv.Interface())

	fmt.Println(reflect.TypeOf(Date{}).NumField())

}

func GetUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}

func SetUnexportedField(field reflect.Value, value interface{}) {
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
		Elem().
		Set(reflect.ValueOf(value))
}
