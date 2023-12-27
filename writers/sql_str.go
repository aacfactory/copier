package writers

import (
	"github.com/modern-go/reflect2"
	"unsafe"
)

type SQLNullStringWriter struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (w *SQLNullStringWriter) Name() string {
	return w.typ.String()
}

func (w *SQLNullStringWriter) Type() reflect2.Type {
	return w.typ
}

func (w *SQLNullStringWriter) ValueType() reflect2.Type {
	return w.valueType.Type()
}

func (w *SQLNullStringWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {

	return
}
