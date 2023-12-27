package writers

import (
	"github.com/modern-go/reflect2"
	"unsafe"
)

type SQLNullTimeWriter struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (w *SQLNullTimeWriter) Name() string {
	return w.typ.String()
}

func (w *SQLNullTimeWriter) Type() reflect2.Type {
	return w.typ
}

func (w *SQLNullTimeWriter) ValueType() reflect2.Type {
	return w.valueType.Type()
}

func (w *SQLNullTimeWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {

	return
}
