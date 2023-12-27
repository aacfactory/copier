package writers

import (
	"github.com/modern-go/reflect2"
	"unsafe"
)

type SQLNullIntWriter struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (w *SQLNullIntWriter) Name() string {
	return w.typ.String()
}

func (w *SQLNullIntWriter) Type() reflect2.Type {
	return w.typ
}

func (w *SQLNullIntWriter) ValueType() reflect2.Type {
	return w.valueType.Type()
}

func (w *SQLNullIntWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {

	return
}
