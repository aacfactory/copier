package writers

import (
	"fmt"
	"github.com/aacfactory/copier/descriptors"
	"github.com/modern-go/reflect2"
	"reflect"
)

type SQLWriter interface {
	Writer
	ValueType() reflect2.Type
}

func IsSQLValue(typ reflect2.Type) bool {
	return typ.Implements(sqlValuerType)
}

func NewUnsafeSQLWriter(typ reflect2.Type) SQLWriter {
	w, err := NewSQLWriter(typ)
	if err != nil {
		panic(err)
		return nil
	}
	return w
}

func NewSQLWriter(typ reflect2.Type) (w SQLWriter, err error) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.(reflect2.PtrType).Elem()
	}
	descriptor := descriptors.DescribeStruct(typ)
	validField := descriptor.Field("Valid")
	if validField == nil {
		err = fmt.Errorf("copier: sql null value writer can not support %s dst type", typ.String())
		return
	}
	if validField.Field[0].Type().Kind() != reflect.Bool {
		err = fmt.Errorf("copier: sql null value writer can not support %s dst type", typ.String())
		return
	}
	// sql
	stringField := descriptor.Field("String")
	if stringField != nil {
		if stringField.Field[0].Type().Kind() != reflect.String {
			err = fmt.Errorf("copier: sql null value writer can not support %s dst type", typ.String())
			return
		}
		w = &SQLNullStringWriter{
			typ:       typ,
			valueType: stringField.Field[0],
			validType: validField.Field[0],
		}
		return
	}
	// bool
	boolField := descriptor.Field("Bool")
	if boolField != nil {
		if boolField.Field[0].Type().Kind() != reflect.Bool {
			err = fmt.Errorf("copier: sql null value writer can not support %s dst type", typ.String())
			return
		}
		w = &SQLNullBoolWriter{
			typ:       typ,
			valueType: boolField.Field[0],
			validType: validField.Field[0],
		}
		return
	}
	// int
	int16Field := descriptor.Field("Int16")
	if int16Field != nil {
		if int16Field.Field[0].Type().Kind() != reflect.Int16 {
			err = fmt.Errorf("copier: sql null value writer can not support %s dst type", typ.String())
			return
		}
		w = &SQLNullIntWriter{
			typ:       typ,
			valueType: int16Field.Field[0],
			validType: validField.Field[0],
		}
		return
	}
	int32Field := descriptor.Field("Int32")
	if int32Field != nil {
		if int32Field.Field[0].Type().Kind() != reflect.Int32 {
			err = fmt.Errorf("copier: sql null value writer can not support %s dst type", typ.String())
			return
		}
		w = &SQLNullIntWriter{
			typ:       typ,
			valueType: int32Field.Field[0],
			validType: validField.Field[0],
		}
		return
	}
	int64Field := descriptor.Field("Int64")
	if int64Field != nil {
		if int64Field.Field[0].Type().Kind() != reflect.Int64 {
			err = fmt.Errorf("copier: sql null value writer can not support %s dst type", typ.String())
			return
		}
		w = &SQLNullIntWriter{
			typ:       typ,
			valueType: int64Field.Field[0],
			validType: validField.Field[0],
		}
		return
	}
	// float
	floatField := descriptor.Field("Float64")
	if floatField != nil {
		if floatField.Field[0].Type().Kind() != reflect.Float64 {
			err = fmt.Errorf("copier: sql null value writer can not support %s dst type", typ.String())
			return
		}
		w = &SQLNullFloatWriter{
			typ:       typ,
			valueType: floatField.Field[0],
			validType: validField.Field[0],
		}
		return
	}
	// byte
	byteField := descriptor.Field("Byte")
	if byteField != nil {
		if byteField.Field[0].Type().Kind() != reflect.Uint8 {
			err = fmt.Errorf("copier: sql null value writer can not support %s dst type", typ.String())
			return
		}
		w = &SQLNullByteWriter{
			typ:       typ,
			valueType: byteField.Field[0],
			validType: validField.Field[0],
		}
		return
	}
	// time
	timeField := descriptor.Field("Time")
	if timeField != nil {
		if timeField.Field[0].Type().RType() == timeType.RType() || timeField.Field[0].Type().Type1().ConvertibleTo(timeType.Type1()) {
			w = &SQLNullTimeWriter{
				typ:       typ,
				valueType: timeField.Field[0],
				validType: validField.Field[0],
			}
			return
		}
		return
	}
	// generic
	var valueType reflect2.StructField
	for _, fieldDescriptor := range descriptor.Fields() {
		if fieldDescriptor.Name == "Valid" {
			continue
		}
		if len(fieldDescriptor.Field) == 1 {
			valueType = fieldDescriptor.Field[0]
			break
		}
	}
	if validField == nil {
		err = fmt.Errorf("copier: sql null value writer can not support %s dst type", typ.String())
		return
	}
	w = &SQLNullGenericWriter{
		typ:       typ,
		valueType: valueType,
		validType: validField.Field[0],
	}
	return
}
