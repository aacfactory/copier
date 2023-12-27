package writers

import (
	"fmt"
	"github.com/modern-go/reflect2"
	"unsafe"
)

func NewMapType(cfg *Writers, typ reflect2.MapType) (v Writer, err error) {
	keyType := typ.Key()
	keyWriter, keyErr := WriterOf(cfg, keyType)
	if keyErr != nil {
		err = fmt.Errorf("copier: not support %s dst type, %v", typ.String(), keyErr)
		return
	}

	valueType := typ.Elem()
	valueWriter, valueErr := WriterOf(cfg, valueType)
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
	typ         reflect2.Type
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

	return
}
