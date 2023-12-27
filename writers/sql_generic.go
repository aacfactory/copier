package writers

import (
	"database/sql/driver"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"unsafe"
)

type SQLNullGenericWriter struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (w *SQLNullGenericWriter) Name() string {
	return w.typ.String()
}

func (w *SQLNullGenericWriter) Type() reflect2.Type {
	return w.typ
}

func (w *SQLNullGenericWriter) ValueType() reflect2.Type {
	return w.valueType.Type()
}

func (w *SQLNullGenericWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {
	if w.typ.RType() == srcType.RType() {
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}
	if srcType.UnsafeIsNil(srcPtr) {
		return
	}
	if w.valueType.Type().RType() == srcType.RType() {
		w.valueType.UnsafeSet(dstPtr, srcPtr)
		return
	}
	if srcType.Type1().ConvertibleTo(w.validType.Type().Type1()) {
		srcPtr = reflect.ValueOf(srcType.PackEFace(srcPtr)).Convert(w.validType.Type().Type1()).UnsafePointer()
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}
	if IsSQLValue(srcType) {
		valuer, isValuer := srcType.PackEFace(srcPtr).(driver.Valuer)
		if !isValuer {
			err = fmt.Errorf("copier: sql null generic writer can not support %s source type", srcType.String())
			return
		}
		value, valueErr := valuer.Value()
		if valueErr != nil {
			err = valueErr
			return
		}
		fmt.Println("sql generic:", value, valueErr)
		if reflect2.IsNil(value) {
			return
		}
		err = w.Write(dstPtr, reflect2.PtrOf(value), reflect2.TypeOf(value))
		return
	}
	err = fmt.Errorf("copier: sql null generic writer can not support %s source type", srcType.String())
	return
}
