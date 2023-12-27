package writers

import (
	"github.com/modern-go/reflect2"
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

	return
}
