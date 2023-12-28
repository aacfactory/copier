package writers

import (
	"database/sql/driver"
	"fmt"
	"github.com/aacfactory/copier/descriptors"
	"github.com/modern-go/reflect2"
	"reflect"
	"unsafe"
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

func (w *StructWriter) Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {
	if srcType.Kind() == reflect.Ptr {
		srcType = srcType.(reflect2.PtrType).Elem()
	}
	// same
	if w.typ.RType() == srcType.RType() {
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}
	// convertable
	if IsConvertible(srcType) {
		srcPtr, srcType = convert(srcPtr, srcType)
	}
	// convertible
	if w.typ.Type1().ConvertibleTo(srcType.Type1()) {
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}
	// sql
	if IsSQLValue(srcType) {
		valuer, isValuer := srcType.PackEFace(srcPtr).(driver.Valuer)
		if !isValuer {
			err = fmt.Errorf("copier: struct writer can not support %s source type", srcType.String())
			return
		}
		value, valueErr := valuer.Value()
		if valueErr != nil {
			err = valueErr
			return
		}
		if reflect2.IsNil(value) {
			return
		}
		err = w.Write(dstPtr, reflect2.PtrOf(value), reflect2.TypeOf(value))
		if err != nil {
			return
		}
		return
	}
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
		dft := field.typ.Type()
		dfp := field.typ.UnsafeGet(dstPtr)
		if dft.UnsafeIsNil(dfp) {
			dfp = dft.UnsafeNew()
		}
		err = field.writer.Write(dfp, sfp, sft)
		if err != nil {
			return
		}
		field.typ.UnsafeSet(dstPtr, dfp)
	}
	return
}

type StructFieldWriter struct {
	typ    reflect2.StructField
	tag    string
	writer Writer
}

type StructFieldWriters []*StructFieldWriter
