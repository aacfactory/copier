package writers

import (
	"github.com/modern-go/reflect2"
	"unsafe"
)

func NewStruct(cfg *Writers, typ reflect2.Type) (w Writer, err error) {

	return
}

type StructWriter struct {
	typ reflect2.Type
}

func (w *StructWriter) Name() string {
	return w.typ.String()
}

func (w *StructWriter) Type() reflect2.Type {
	return w.typ
}

func (w *StructWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {

	return
}
