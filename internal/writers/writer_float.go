package writers

import (
	"database/sql"
	"fmt"
	"github.com/aacfactory/copier/internal/commons"
	"github.com/modern-go/reflect2"
	"reflect"
	"strconv"
	"unsafe"
)

func NewFloatWriter() Writer {
	return &FloatWriter{
		typ: reflect2.TypeOf(float64(0)),
	}
}

type FloatWriter struct {
	typ reflect2.Type
}

func (w *FloatWriter) Write(dstPtr unsafe.Pointer, reader Reader) (err error) {
	srcType := reader.Type()
	srcPtr := reader.Read()
	switch srcType.Kind() {
	case reflect.String:
		s := *(*string)(srcPtr)
		n, nErr := strconv.ParseFloat(s, 64)
		if nErr != nil {
			err = fmt.Errorf("copier: float writer can not support %s type reader, src value is not float format string", srcType.String())
			return
		}
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(n))
		break
	case reflect.Bool:
		b := *(*bool)(srcPtr)
		if b {
			w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(float64(1)))
		}
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n := *(*int64)(srcPtr)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(float64(n)))
		break
	case reflect.Float32, reflect.Float64:
		w.typ.UnsafeSet(dstPtr, srcPtr)
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u := *(*uint64)(srcPtr)
		w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(float64(u)))
		break
	default:
		// sql: null string
		if srcType.Type1().ConvertibleTo(sqlNullStringType.Type1()) {
			nv := new(sql.NullString)
			nptr := reflect2.PtrOf(nv)
			sqlNullStringType.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				n, nErr := strconv.ParseFloat(nv.String, 64)
				if nErr != nil {
					err = fmt.Errorf("copier: float writer can not support %s type reader, src value is not float format string", srcType.String())
					return
				}
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(n))
			}
			break
		}
		// sql: null bool
		if srcType.Type1().ConvertibleTo(sqlNullBoolType.Type1()) {
			nv := new(sql.NullBool)
			nptr := reflect2.PtrOf(nv)
			sqlNullBoolType.UnsafeSet(nptr, srcPtr)
			if nv.Valid && nv.Bool {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(float64(1)))
			}
			break
		}
		// sql: null int16
		if srcType.Type1().ConvertibleTo(sqlNullInt16Type.Type1()) {
			nv := new(sql.NullInt16)
			nptr := reflect2.PtrOf(nv)
			sqlNullInt16Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(float64(nv.Int16)))
			}
			break
		}
		// sql: null int32
		if srcType.Type1().ConvertibleTo(sqlNullInt32Type.Type1()) {
			nv := new(sql.NullInt32)
			nptr := reflect2.PtrOf(nv)
			sqlNullInt32Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(float64(nv.Int32)))
			}
			break
		}
		// sql: null int64
		if srcType.Type1().ConvertibleTo(sqlNullInt64Type.Type1()) {
			nv := new(sql.NullInt64)
			nptr := reflect2.PtrOf(nv)
			sqlNullInt64Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(float64(nv.Int64)))
			}
			break
		}
		// sql: null float
		if srcType.Type1().ConvertibleTo(sqlNullFloat64Type.Type1()) {
			nv := new(sql.NullFloat64)
			nptr := reflect2.PtrOf(nv)
			sqlNullFloat64Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.typ.UnsafeSet(dstPtr, reflect2.PtrOf(nv.Float64))
			}
			break
		}
		err = fmt.Errorf("copier: float writer can not support %s type reader", srcType.String())
		return
	}
	w.typ.UnsafeSet(dstPtr, srcPtr)
	return
}

func NewNullFloatWriter(typ reflect2.Type) Writer {
	if typ.Kind() == reflect.Ptr {
		typ = typ.(reflect2.PtrType).Elem()
	}
	descriptor := commons.DescribeStruct("", typ)
	var valueType reflect2.StructField
	var validType reflect2.StructField
	for _, field := range descriptor.Fields {
		for _, f := range field.Field {
			if f.Name() == "Float64" && f.Type().Kind() == reflect.Float64 {
				valueType = f.Type().(reflect2.StructField)
				continue
			}
			if f.Name() == "Valid" && f.Type().Kind() == reflect.Bool {
				validType = f.Type().(reflect2.StructField)
				continue
			}
		}
	}
	return &NullFloatWriter{
		typ:       typ,
		valueType: valueType,
		validType: validType,
	}
}

type NullFloatWriter struct {
	typ       reflect2.Type
	valueType reflect2.StructField
	validType reflect2.StructField
}

func (w *NullFloatWriter) Write(dstPtr unsafe.Pointer, reader Reader) (err error) {
	srcType := reader.Type()
	srcPtr := reader.Read()
	if w.typ.RType() == srcType.RType() {
		w.typ.UnsafeSet(dstPtr, srcPtr)
		return
	}
	switch srcType.Kind() {
	case reflect.String:
		s := *(*string)(srcPtr)
		n, nErr := strconv.ParseFloat(s, 64)
		if nErr != nil {
			err = fmt.Errorf("copier: null float writer can not support %s type reader, src value is not float format string", srcType.String())
			return
		}
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(n))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Bool:
		b := *(*bool)(srcPtr)
		if b {
			w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(float64(1)))
		} else {
			w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(float64(0)))
		}
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n := *(*int64)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(float64(n)))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Float32, reflect.Float64:
		n := *(*float64)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(n))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n := *(*uint64)(srcPtr)
		w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(float64(n)))
		w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
		break
	default:
		// sql: null string
		if srcType.Type1().ConvertibleTo(sqlNullStringType.Type1()) {
			nv := new(sql.NullString)
			nptr := reflect2.PtrOf(nv)
			sqlNullStringType.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				n, nErr := strconv.ParseFloat(nv.String, 64)
				if nErr != nil {
					err = fmt.Errorf("copier: null float writer can not support %s type reader, src value is not float format string", srcType.String())
					return
				}
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(n))
				w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		// sql: null bool
		if srcType.Type1().ConvertibleTo(sqlNullBoolType.Type1()) {
			nv := new(sql.NullBool)
			nptr := reflect2.PtrOf(nv)
			sqlNullBoolType.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				if nv.Bool {
					w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(float64(1)))
				}
				w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		// sql: null int16
		if srcType.Type1().ConvertibleTo(sqlNullInt16Type.Type1()) {
			nv := new(sql.NullInt16)
			nptr := reflect2.PtrOf(nv)
			sqlNullInt16Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(float64(nv.Int16)))
				w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		// sql: null int32
		if srcType.Type1().ConvertibleTo(sqlNullInt32Type.Type1()) {
			nv := new(sql.NullInt32)
			nptr := reflect2.PtrOf(nv)
			sqlNullInt32Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(float64(nv.Int32)))
				w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		// sql: null int64
		if srcType.Type1().ConvertibleTo(sqlNullInt64Type.Type1()) {
			nv := new(sql.NullInt64)
			nptr := reflect2.PtrOf(nv)
			sqlNullInt64Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(float64(nv.Int64)))
				w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		// sql: null float
		if srcType.Type1().ConvertibleTo(sqlNullFloat64Type.Type1()) {
			nv := new(sql.NullFloat64)
			nptr := reflect2.PtrOf(nv)
			sqlNullFloat64Type.UnsafeSet(nptr, srcPtr)
			if nv.Valid {
				w.valueType.UnsafeSet(dstPtr, reflect2.PtrOf(nv.Float64))
				w.validType.UnsafeSet(dstPtr, reflect2.PtrOf(true))
			}
			break
		}
		err = fmt.Errorf("copier: null float writer can not support %s type reader", srcType.String())
		return
	}
	return
}
