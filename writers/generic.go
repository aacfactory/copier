package writers

import (
	"errors"
	"fmt"
	"github.com/modern-go/reflect2"
	"unsafe"
)

type GenericWriter interface {
	UnsafeWrite(srcPtr unsafe.Pointer, srcType reflect2.Type) (err error)
}

func NewGenericWriter(typ reflect2.Type) *GenericInterfaceWriter {
	return &GenericInterfaceWriter{
		typ: typ,
	}
}

type GenericInterfaceWriter struct {
	typ reflect2.Type
}

func (w *GenericInterfaceWriter) Name() string {
	return w.typ.String()
}

func (w *GenericInterfaceWriter) Type() reflect2.Type {
	return w.typ
}

func (w *GenericInterfaceWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {
	// convertable
	if IsConvertible(srcType) {
		srcPtr, srcType = convert(srcPtr, srcType)
	}
	uw, ok := w.typ.PackEFace(dstPtr).(GenericWriter)
	if !ok {
		err = fmt.Errorf("copier: generic writer can not support %s dst type", reflect2.TypeOf(w.typ.UnsafeIndirect(dstPtr)))
		return
	}
	if uw == nil {
		return
	}
	err = uw.UnsafeWrite(srcPtr, srcType)
	if err != nil {
		err = errors.Join(
			fmt.Errorf("copier: %s generic writer failed", w.Name()),
			err,
		)
		return
	}
	return
}
