package writers

import (
	"errors"
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
)

func NewMapWriter(cfg *Writers, typ reflect2.MapType) (v Writer, err error) {
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

func (w *MapWriter) Write(dst any, src any) (err error) {
	if src == nil {
		return
	}
	srcType := reflect2.TypeOfPtr(src).Elem()
	srcPtr := reflect2.PtrOf(src)
	dstPtr := reflect2.PtrOf(dst)

	fmt.Println(reflect.TypeOf(src))

	if w.typ.RType() == srcType.RType() {
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}

	switch srcType.Kind() {
	case reflect.Map:
		smt := srcType.(reflect2.MapType)
		fmt.Println(src)
		fmt.Println(srcType.Indirect(src))
		iter := smt.Iterate(src)
		for {
			if !iter.HasNext() {
				break
			}
			sk, sv := iter.Next()
			dk := w.keyType.New()
			kErr := w.keyWriter.Write(dk, sk)
			if kErr != nil {
				err = errors.Join(
					fmt.Errorf("copier: %s map writer write key faield", w.typ.String()),
					kErr,
				)
				return
			}
			dv := w.valueType.New()
			vErr := w.valueWriter.Write(dv, sv)
			if vErr != nil {
				err = errors.Join(
					fmt.Errorf("copier: %s map writer write value faield", w.typ.String()),
					vErr,
				)
				return
			}
			w.typ.SetIndex(dst, dk, dv)
		}
		break
	case reflect.Struct:
		// convertable
		if convertible, ok := src.(Convertible); ok {
			value := convertible.Convert()
			if value == nil {
				return
			}
			err = w.Write(dst, reflect2.TypeOf(value).PackEFace(reflect2.PtrOf(value)))
			return
		}
		err = fmt.Errorf("copier: map writer can not support %s source type", srcType.String())
		return
	case reflect.Ptr:
		// convertable
		if convertible, ok := src.(Convertible); ok {
			value := convertible.Convert()
			if value == nil {
				return
			}
			err = w.Write(dst, reflect2.TypeOf(value).PackEFace(reflect2.PtrOf(value)))
			return
		}
		err = w.Write(dst, srcType.Indirect(src))
		break
	default:
		err = fmt.Errorf("copier: map writer can not support %s source type", srcType.String())
		return
	}
	return
}
