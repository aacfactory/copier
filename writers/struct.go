package writers

import (
	"database/sql/driver"
	"fmt"
	"github.com/aacfactory/copier/descriptors"
	"github.com/modern-go/reflect2"
)

func NewStruct(cfg *Writers, typ reflect2.Type) (w Writer, err error) {
	// sql
	if typ.Implements(sqlValuerType) {
		w, err = NewSQLWriter(typ)
		return
	}
	// time
	if typ.RType() == timeType.RType() || typ.Type1().ConvertibleTo(timeType.Type1()) {
		w = NewTimeWriter(typ)
		return
	}
	// set processing
	sw := &StructWriter{
		tagKey: cfg.tagKey,
		typ:    typ,
	}
	cfg.addProcessing(typ, sw)
	// fields
	fields := make(StructFieldWriters, 0, 1)
	desc := descriptors.DescribeStruct(typ)
	for _, fieldDescriptor := range desc.Fields() {
		sft := fieldDescriptor.StructField()
		ft := sft.Type()
		fw, fwErr := cfg.Get(ft)
		if fwErr != nil {
			err = fwErr
			return
		}
		tag, _ := sft.Tag().Lookup(cfg.tagKey)
		fields = append(fields, &StructFieldWriter{
			typ:    sft,
			tag:    tag,
			writer: fw,
		})
	}
	sw.fields = fields
	// rem processing
	cfg.removeProcessing(typ)
	// return
	w = sw
	return
}

type StructWriter struct {
	tagKey string
	typ    reflect2.Type
	fields StructFieldWriters
}

func (w *StructWriter) Name() string {
	return w.typ.String()
}

func (w *StructWriter) Type() reflect2.Type {
	return w.typ
}

func (w *StructWriter) Write(dst any, src any) (err error) {
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

	// convertible
	if w.typ.Type1().ConvertibleTo(srcType.Type1()) {
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}
	// sql
	if valuer, ok := src.(driver.Valuer); ok {
		value, valueErr := valuer.Value()
		if valueErr != nil {
			err = valueErr
			return
		}
		err = w.Write(dst, reflect2.TypeOf(value).PackEFace(reflect2.PtrOf(value)))
		return
	}
	fmt.Println(srcType, reflect2.TypeOfPtr(src))
	desc := descriptors.DescribeStruct(srcType)
	for _, field := range w.fields {
		var srcField *descriptors.StructFieldDescriptor
		if field.tag == "" {
			srcField = desc.Field(field.typ.Name())
		} else {
			srcField = desc.FieldByTag(w.tagKey, field.tag)
		}
		if srcField == nil {
			continue
		}
		sfp, sft, sfErr := srcField.ValueOf(srcPtr)
		if sfErr != nil {
			err = sfErr
			return
		}
		if sft.UnsafeIsNil(sfp) {
			continue
		}
		sf := sft.PackEFace(sfp)
		dft := field.typ.Type()
		df := field.typ.Get(dst)
		if dft.IsNil(df) {
			df = dft.New()
		}
		err = field.writer.Write(df, sf)
		if err != nil {
			return
		}
		field.typ.Set(dst, df)
	}
	return
}

type StructFieldWriter struct {
	typ    reflect2.StructField
	tag    string
	writer Writer
}

type StructFieldWriters []*StructFieldWriter
