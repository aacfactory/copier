package descriptors

import (
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"unsafe"
)

type SQLDescriptor interface {
	Descriptor
	ValueOf(ptr unsafe.Pointer) (v unsafe.Pointer, valid bool)
}

// +-------------------------------------------------------------------------------------------------------------------+

func DescribeNullString(typ reflect2.Type) (v SQLDescriptor, err error) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.(reflect2.PtrType).Elem()
	}
	descriptor := DescribeStruct("", typ)
	var valueType reflect2.StructField
	var validType reflect2.StructField
	for _, field := range descriptor.Fields() {
		for _, f := range field.Field {
			if f.Name() == "String" && f.Type().Kind() == reflect.String {
				valueType = f
				continue
			}
			if f.Name() == "Valid" && f.Type().Kind() == reflect.Bool {
				validType = f
				continue
			}
		}
	}
	if valueType == nil || validType == nil {
		err = fmt.Errorf("copier: null string descriptor can not support %s type reader", typ.String())
		return
	}
	v = &NullStringDescriptor{
		typ:       typ,
		valueType: valueType,
		validType: validType,
	}
	return
}

type NullStringDescriptor struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (desc *NullStringDescriptor) Type() reflect2.Type {
	return desc.typ
}

func (desc *NullStringDescriptor) ValueOf(ptr unsafe.Pointer) (v unsafe.Pointer, valid bool) {
	validPtr := desc.validType.UnsafeGet(ptr)
	valid = *(*bool)(validPtr)
	if valid {
		v = desc.valueType.UnsafeGet(ptr)
	}
	return
}

// +-------------------------------------------------------------------------------------------------------------------+

func DescribeNullBool(typ reflect2.Type) (v SQLDescriptor, err error) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.(reflect2.PtrType).Elem()
	}
	descriptor := DescribeStruct("", typ)
	var valueType reflect2.StructField
	var validType reflect2.StructField
	for _, field := range descriptor.Fields() {
		for _, f := range field.Field {
			if f.Name() == "Bool" && f.Type().Kind() == reflect.Bool {
				valueType = f
				continue
			}
			if f.Name() == "Valid" && f.Type().Kind() == reflect.Bool {
				validType = f
				continue
			}
		}
	}
	if valueType == nil || validType == nil {
		err = fmt.Errorf("copier: null bool descriptor can not support %s type reader", typ.String())
		return
	}
	v = &NullBoolDescriptor{
		typ:       typ,
		valueType: valueType,
		validType: validType,
	}
	return
}

type NullBoolDescriptor struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (desc *NullBoolDescriptor) Type() reflect2.Type {
	return desc.typ
}

func (desc *NullBoolDescriptor) ValueOf(ptr unsafe.Pointer) (v unsafe.Pointer, valid bool) {
	validPtr := desc.validType.UnsafeGet(ptr)
	valid = *(*bool)(validPtr)
	if valid {
		v = desc.valueType.UnsafeGet(ptr)
	}
	return
}

// +-------------------------------------------------------------------------------------------------------------------+

func DescribeNullInt(typ reflect2.Type) (v SQLDescriptor, err error) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.(reflect2.PtrType).Elem()
	}
	descriptor := DescribeStruct("", typ)
	var valueType reflect2.StructField
	var validType reflect2.StructField
	for _, field := range descriptor.Fields() {
		for _, f := range field.Field {
			if f.Name() == "Int16" && f.Type().Kind() == reflect.Int16 {
				valueType = f
				continue
			}
			if f.Name() == "Int32" && f.Type().Kind() == reflect.Int32 {
				valueType = f
				continue
			}
			if f.Name() == "Int64" && f.Type().Kind() == reflect.Int64 {
				valueType = f
				continue
			}
			if f.Name() == "Valid" && f.Type().Kind() == reflect.Bool {
				validType = f
				continue
			}
		}
	}
	if valueType == nil || validType == nil {
		err = fmt.Errorf("copier: null int descriptor can not support %s type reader", typ.String())
		return
	}
	v = &NullIntDescriptor{
		typ:       typ,
		valueType: valueType,
		validType: validType,
	}
	return
}

type NullIntDescriptor struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (desc *NullIntDescriptor) Type() reflect2.Type {
	return desc.typ
}

func (desc *NullIntDescriptor) ValueOf(ptr unsafe.Pointer) (v unsafe.Pointer, valid bool) {
	validPtr := desc.validType.UnsafeGet(ptr)
	valid = *(*bool)(validPtr)
	if valid {
		v = desc.valueType.UnsafeGet(ptr)
	}
	return
}

// +-------------------------------------------------------------------------------------------------------------------+

func DescribeNullFloat(typ reflect2.Type) (v SQLDescriptor, err error) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.(reflect2.PtrType).Elem()
	}
	descriptor := DescribeStruct("", typ)
	var valueType reflect2.StructField
	var validType reflect2.StructField
	for _, field := range descriptor.Fields() {
		for _, f := range field.Field {
			if f.Name() == "Float64" && f.Type().Kind() == reflect.Float64 {
				valueType = f
				continue
			}
			if f.Name() == "Valid" && f.Type().Kind() == reflect.Bool {
				validType = f
				continue
			}
		}
	}
	if valueType == nil || validType == nil {
		err = fmt.Errorf("copier: null float descriptor can not support %s type reader", typ.String())
		return
	}
	v = &NullFloatDescriptor{
		typ:       typ,
		valueType: valueType,
		validType: validType,
	}
	return
}

type NullFloatDescriptor struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (desc *NullFloatDescriptor) Type() reflect2.Type {
	return desc.typ
}

func (desc *NullFloatDescriptor) ValueOf(ptr unsafe.Pointer) (v unsafe.Pointer, valid bool) {
	validPtr := desc.validType.UnsafeGet(ptr)
	valid = *(*bool)(validPtr)
	if valid {
		v = desc.valueType.UnsafeGet(ptr)
	}
	return
}

// +-------------------------------------------------------------------------------------------------------------------+

func DescribeNullByte(typ reflect2.Type) (v SQLDescriptor, err error) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.(reflect2.PtrType).Elem()
	}
	descriptor := DescribeStruct("", typ)
	var valueType reflect2.StructField
	var validType reflect2.StructField
	for _, field := range descriptor.Fields() {
		for _, f := range field.Field {
			if f.Name() == "Byte" && f.Type().Kind() == reflect.Uint8 {
				valueType = f
				continue
			}
			if f.Name() == "Valid" && f.Type().Kind() == reflect.Bool {
				validType = f
				continue
			}
		}
	}
	if valueType == nil || validType == nil {
		err = fmt.Errorf("copier: null byte descriptor can not support %s type reader", typ.String())
		return
	}
	v = &NullByteDescriptor{
		typ:       typ,
		valueType: valueType,
		validType: validType,
	}
	return
}

type NullByteDescriptor struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (desc *NullByteDescriptor) Type() reflect2.Type {
	return desc.typ
}

func (desc *NullByteDescriptor) ValueOf(ptr unsafe.Pointer) (v unsafe.Pointer, valid bool) {
	validPtr := desc.validType.UnsafeGet(ptr)
	valid = *(*bool)(validPtr)
	if valid {
		v = desc.valueType.UnsafeGet(ptr)
	}
	return
}

// +-------------------------------------------------------------------------------------------------------------------+

func DescribeNullTime(typ reflect2.Type) (v SQLDescriptor, err error) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.(reflect2.PtrType).Elem()
	}
	descriptor := DescribeStruct("", typ)
	var valueType reflect2.StructField
	var validType reflect2.StructField
	for _, field := range descriptor.Fields() {
		for _, f := range field.Field {
			if f.Name() == "Time" && f.Type().Type1().ConvertibleTo(timeType.Type1()) {
				valueType = f
				continue
			}
			if f.Name() == "Valid" && f.Type().Kind() == reflect.Bool {
				validType = f
				continue
			}
		}
	}
	if valueType == nil || validType == nil {
		err = fmt.Errorf("copier: null time descriptor can not support %s type reader", typ.String())
		return
	}
	v = &NullTimeDescriptor{
		typ:       typ,
		valueType: valueType,
		validType: validType,
	}
	return
}

type NullTimeDescriptor struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (desc *NullTimeDescriptor) Type() reflect2.Type {
	return desc.typ
}

func (desc *NullTimeDescriptor) ValueOf(ptr unsafe.Pointer) (v unsafe.Pointer, valid bool) {
	validPtr := desc.validType.UnsafeGet(ptr)
	valid = *(*bool)(validPtr)
	if valid {
		v = desc.valueType.UnsafeGet(ptr)
	}
	return
}
