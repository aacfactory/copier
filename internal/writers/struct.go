package writers

import (
	"github.com/aacfactory/copier/internal/commons"
	"unsafe"
)

func NewStructWriter() Writer {

	return nil
}

type StructWriter struct {
	descriptor *commons.StructDescriptor
}

func (w *StructWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer) {

	return
}
