package writers

import (
	"github.com/modern-go/reflect2"
	"unsafe"
)

type SQLNullBoolWriter struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (w *SQLNullBoolWriter) Name() string {
	return w.typ.String()
}

func (w *SQLNullBoolWriter) Type() reflect2.Type {
	return w.typ
}

func (w *SQLNullBoolWriter) ValueType() reflect2.Type {
	return w.valueType.Type()
}

func (w *SQLNullBoolWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {

	return
}
