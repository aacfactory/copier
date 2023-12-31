package writers

import (
	"errors"
	"fmt"
	"github.com/modern-go/reflect2"
)

type GenericWriter interface {
	UnsafeWrite(obj any) (err error)
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

func (w *GenericInterfaceWriter) Write(dst any, src any) (err error) {
	if src == nil {
		return
	}
	srcType := reflect2.TypeOfPtr(src).Elem()
	srcPtr := reflect2.PtrOf(src)
	dstPtr := reflect2.PtrOf(dst)

	if w.typ.RType() == srcType.RType() {
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}

	// convertable
	if convertible, ok := src.(Convertible); ok {
		value := convertible.Convert()
		if value == nil {
			return
		}
		err = w.Write(dst, reflect2.TypeOf(value).PackEFace(reflect2.PtrOf(value)))
		return
	}

	uw, ok := dst.(GenericWriter)
	if !ok {
		err = fmt.Errorf("copier: generic writer can not support %s dst type", reflect2.TypeOf(dst))
		return
	}
	if uw == nil {
		return
	}

	err = uw.UnsafeWrite(srcType.Indirect(src))
	if err != nil {
		err = errors.Join(
			fmt.Errorf("copier: %s generic writer failed", w.Name()),
			err,
		)
		return
	}
	return
}
