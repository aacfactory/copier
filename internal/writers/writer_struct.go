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

func (w *StructWriter) Write(dstPtr unsafe.Pointer, reader Reader) (err error) {
	srcType := reader.Type()

	return
}
