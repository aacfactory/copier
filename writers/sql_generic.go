package writers

import (
	"database/sql/driver"
	"fmt"
	"github.com/modern-go/reflect2"
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

func (w *SQLNullGenericWriter) Write(dst any, src any) (err error) {
	if src == nil {
		return
	}
	srcType := reflect2.TypeOfPtr(src).Elem()
	srcPtr := reflect2.PtrOf(src)
	dstPtr := reflect2.PtrOf(dst)

	if w.typ.RType() == srcType.RType() {
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}

	if w.valueType.Type().RType() == srcType.RType() {
		w.valueType.UnsafeSet(dstPtr, srcPtr)
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		return
	}

	// convertable
	if convertible, ok := src.(Convertible); ok {
		value := convertible.Convert()
		if value == nil {
			return
		}
		err = w.Write(dst, reflect2.TypeOf(value).PackEFace(reflect2.PtrOf(value)))
		return
	}
	// sql
	if valuer, ok := src.(driver.Valuer); ok {
		value, valueErr := valuer.Value()
		if valueErr != nil {
			err = valueErr
			return
		}
		if value == nil {
			return
		}
		err = w.Write(dst, reflect2.TypeOf(value).PackEFace(reflect2.PtrOf(value)))
		return
	}

	err = fmt.Errorf("copier: sql null generic writer can not support %s source type", srcType.String())
	return
}
