package writers

import (
	"github.com/modern-go/reflect2"
	"unsafe"
)

type SQLNullByteWriter struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (w *SQLNullByteWriter) Name() string {
	return w.typ.String()
}

func (w *SQLNullByteWriter) Type() reflect2.Type {
	return w.typ
}

func (w *SQLNullByteWriter) ValueType() reflect2.Type {
	return w.valueType.Type()
}

func (w *SQLNullByteWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {

	return
}
