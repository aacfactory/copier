package writers

import (
	"github.com/aacfactory/copier/internal/commons"
	"github.com/modern-go/reflect2"
	"unsafe"
)

type StructReader struct {
	descriptor *commons.StructDescriptor
	fields     map[string]Reader
	ptr        unsafe.Pointer
}

func (r *StructReader) Type() (typ reflect2.Type) {
	typ = r.descriptor.Type
	return
}

func (r *StructReader) Read() (v unsafe.Pointer) {
	v = r.ptr
	return
}

func (r *StructReader) Field(name string) (fieldReader Reader) {

	return
}
