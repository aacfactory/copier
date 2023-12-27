package writers

import (
	"github.com/modern-go/reflect2"
	"unsafe"
)

type SQLNullFloatWriter struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (w *SQLNullFloatWriter) Name() string {
	return w.typ.String()
}

func (w *SQLNullFloatWriter) Type() reflect2.Type {
	return w.typ
}

func (w *SQLNullFloatWriter) ValueType() reflect2.Type {
	return w.valueType.Type()
}

func (w *SQLNullFloatWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {

	return
}
