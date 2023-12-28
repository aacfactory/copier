package writers

import (
	"fmt"
	"github.com/modern-go/reflect2"
	"unsafe"
)

func NewPtrWriter(cfg *Writers, typ reflect2.PtrType) (w Writer, err error) {
	// generic writer
	if typ.Implements(unsafeWriterType) {
		w = NewGenericWriter(typ)
		return
	}
	// sql
	if typ.Implements(sqlValuerType) {
		w, err = NewSQLWriter(typ)
		return
	}
	// text
	if typ.Implements(textUnmarshalerType) {
		w = NewTextUnmarshalerWriter(typ)
		return
	}
	// ptr
	elemType := typ.Elem()
	elemWriter, elemErr := cfg.Get(elemType)
	if elemErr != nil {
		err = fmt.Errorf("copier: ptr not support %s dst type, %v", typ.String(), elemErr)
		return
	}
	w = &PtrWriter{
		typ:        typ,
		elemType:   elemType,
		elemWriter: elemWriter,
	}
	return
}

type PtrWriter struct {
	typ        reflect2.PtrType
	elemType   reflect2.Type
	elemWriter Writer
}

func (w *PtrWriter) Name() string {
	return w.typ.String()
}

func (w *PtrWriter) Type() reflect2.Type {
	return w.typ
}

func (w *PtrWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {
	if srcType.UnsafeIsNil(srcPtr) {
		return
	}
	err = w.elemWriter.Write(dstPtr, srcPtr, srcType)
	return
}
