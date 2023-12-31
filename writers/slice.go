package writers

import (
	"errors"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
)

func NewSliceType(cfg *Writers, typ reflect2.SliceType) (v Writer, err error) {
	elemType := typ.Elem()
	if elemType.Kind() == reflect.Uint8 {
		v = NewBytesWriter(typ)
		return
	}
	elemWriter, elemErr := cfg.Get(elemType)
	if elemErr != nil {
		err = fmt.Errorf("copier: slice writer not support %s dst type, %v", typ.String(), elemErr)
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

func (w *SliceWriter) Write(dst any, src any) (err error) {
	if src == nil {
		return
	}
	srcType := reflect2.TypeOfPtr(src).Elem()
	srcPtr := reflect2.PtrOf(src)
	dstPtr := reflect2.PtrOf(dst)

	if srcType.Kind() != reflect.Slice {
		err = fmt.Errorf("copier: slice writer can not support %s source type", srcType.String())
		return
	}

	if w.typ.RType() == srcType.RType() {
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}
	sst := srcType.(reflect2.SliceType)
	srcLen := sst.UnsafeLengthOf(srcPtr)
	if srcLen == 0 {
		return
	}
	for i := 0; i < srcLen; i++ {
		se := sst.GetIndex(src, i)
		de := w.elemType.New()
		elemWErr := w.elemWriter.Write(de, se)
		if elemWErr != nil {
			err = errors.Join(
				fmt.Errorf("copier: %s slice writer write element faield", w.typ.String()),
				elemWErr,
			)
			return
		}
		w.typ.Append(dst, de)
	}
	return
}
