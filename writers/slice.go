package writers

import (
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"unsafe"
)

func NewSliceType(cfg *Writers, typ reflect2.SliceType) (v Writer, err error) {
	elemType := typ.Elem()
	if elemType.Kind() == reflect.Uint8 {
		v = NewBytesWriter(typ)
		return
	}
	elemWriter, elemErr := cfg.Get(elemType)
	if elemErr != nil {
		err = fmt.Errorf("copier: not support %s dst type, %v", typ.String(), elemErr)
		return
	}
	v = &SliceWriter{
		typ:        typ,
		elemType:   elemType,
		elemWriter: elemWriter,
	}
	return
}

type SliceWriter struct {
	typ        reflect2.SliceType
	elemType   reflect2.Type
	elemWriter Writer
}

func (w *SliceWriter) Name() string {
	return w.typ.String()
}

func (w *SliceWriter) Type() reflect2.Type {
	return w.typ
}

func (w *SliceWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {
	//TODO implement me
	panic("implement me")
}
