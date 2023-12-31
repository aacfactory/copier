package writers

import (
	"fmt"
	"github.com/modern-go/reflect2"
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

func (w *PtrWriter) Write(dst any, src any) (err error) {
	if src == nil {
		return
	}

	srcPtrType := reflect2.TypeOfPtr(src)
	srcType := srcPtrType.Elem()
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

	err = w.elemWriter.Write(dst, srcType.PackEFace(reflect2.PtrOf(srcPtrType.Indirect(src))))
	return
}
