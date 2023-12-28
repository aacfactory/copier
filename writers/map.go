package writers

import (
	"errors"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"unsafe"
)

func NewMapType(cfg *Writers, typ reflect2.MapType) (v Writer, err error) {
	keyType := typ.Key()
	keyWriter, keyErr := cfg.Get(keyType)
	if keyErr != nil {
		err = fmt.Errorf("copier: not support %s dst type, %v", typ.String(), keyErr)
		return
	}

	valueType := typ.Elem()
	valueWriter, valueErr := cfg.Get(valueType)
	if valueErr != nil {
		err = fmt.Errorf("copier: not support %s dst type, %v", typ.String(), valueErr)
		return
	}
	v = &MapWriter{
		typ:         typ,
		keyType:     keyType,
		valueType:   valueType,
		keyWriter:   keyWriter,
		valueWriter: valueWriter,
	}
	return
}

type MapWriter struct {
	typ         reflect2.MapType
	keyType     reflect2.Type
	valueType   reflect2.Type
	keyWriter   Writer
	valueWriter Writer
}

func (w *MapWriter) Name() string {
	return w.typ.String()
}

func (w *MapWriter) Type() reflect2.Type {
	return w.typ
}

func (w *MapWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {
	if srcType.UnsafeIsNil(srcPtr) {
		return
	}
	if w.typ.RType() == srcType.RType() {
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}
	// convertable
	if IsConvertible(srcType) {
		srcPtr, srcType = convert(srcPtr, srcType)
	}
	if srcType.Kind() != reflect.Map {
		err = fmt.Errorf("copier: map writer can not support %s source type", srcType.String())
		return
	}
	smt := srcType.(reflect2.MapType)
	smkt := smt.Key()
	smvt := smt.Elem()
	iterator := smt.UnsafeIterate(srcPtr)
	for {
		if !iterator.HasNext() {
			break
		}
		skp, sep := iterator.UnsafeNext()
		dkp := w.keyType.UnsafeNew()
		kErr := w.keyWriter.Write(dkp, skp, smkt)
		if kErr != nil {
			err = errors.Join(
				fmt.Errorf("copier: %s map writer write key faield", w.typ.String()),
				kErr,
			)
			return
		}
		dvp := w.valueType.UnsafeNew()
		vErr := w.valueWriter.Write(dvp, sep, smvt)
		if vErr != nil {
			err = errors.Join(
				fmt.Errorf("copier: %s map writer write value faield", w.typ.String()),
				vErr,
			)
			return
		}
		w.typ.UnsafeSetIndex(dstPtr, dkp, dvp)
	}
	return
}
