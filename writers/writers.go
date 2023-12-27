package writers

import (
	"database/sql"
	"encoding/binary"
	"fmt"
	"github.com/modern-go/reflect2"
	"golang.org/x/sync/singleflight"
	"reflect"
	"sync"
	"unsafe"
)

type Writer interface {
	Name() string
	Type() reflect2.Type
	Write(dstPtr unsafe.Pointer, srcPtr unsafe.Pointer, srcType reflect2.Type) (err error)
}

func New(tagKey string) *Writers {
	w := &Writers{
		tagKey:       tagKey,
		cache:        sync.Map{},
		processing:   sync.Map{},
		group:        new(singleflight.Group),
		groupKeyPool: new(sync.Pool),
	}
	// native
	w.set(reflect2.RTypeOf(""), NewStringWriter())
	w.set(reflect2.RTypeOf(false), NewBoolWriter())
	w.set(reflect2.RTypeOf(0), NewIntWriter())
	w.set(reflect2.RTypeOf(int8(0)), NewIntWriter())
	w.set(reflect2.RTypeOf(int16(0)), NewIntWriter())
	w.set(reflect2.RTypeOf(int32(0)), NewIntWriter())
	w.set(reflect2.RTypeOf(int64(0)), NewIntWriter())
	w.set(reflect2.RTypeOf(float32(0)), NewFloatWriter())
	w.set(reflect2.RTypeOf(float64(0)), NewFloatWriter())
	w.set(reflect2.RTypeOf(uint(0)), NewUintWriter())
	w.set(reflect2.RTypeOf(uint8(0)), NewUintWriter())
	w.set(reflect2.RTypeOf(uint16(0)), NewUintWriter())
	w.set(reflect2.RTypeOf(uint32(0)), NewUintWriter())
	w.set(reflect2.RTypeOf(uint64(0)), NewUintWriter())
	w.set(reflect2.RTypeOf([]byte{}), NewBytesWriter(bytesType))
	// sql
	w.set(reflect2.RTypeOf(sql.NullString{}), NewUnsafeSQLWriter(sqlNullStringType))
	w.set(reflect2.RTypeOf(sql.NullByte{}), NewUnsafeSQLWriter(sqlNullByteType))
	w.set(reflect2.RTypeOf(sql.NullBool{}), NewUnsafeSQLWriter(sqlNullBoolType))
	w.set(reflect2.RTypeOf(sql.NullInt16{}), NewUnsafeSQLWriter(sqlNullInt16Type))
	w.set(reflect2.RTypeOf(sql.NullInt32{}), NewUnsafeSQLWriter(sqlNullInt32Type))
	w.set(reflect2.RTypeOf(sql.NullInt64{}), NewUnsafeSQLWriter(sqlNullInt64Type))
	w.set(reflect2.RTypeOf(sql.NullFloat64{}), NewUnsafeSQLWriter(sqlNullFloat64Type))
	w.set(reflect2.RTypeOf(sql.NullTime{}), NewUnsafeSQLWriter(sqlNullTimeType))
	return w
}

type Writers struct {
	tagKey       string
	cache        sync.Map
	processing   sync.Map
	group        *singleflight.Group
	groupKeyPool *sync.Pool
}

func (cfg *Writers) set(rtype uintptr, w Writer) {
	cfg.cache.Store(rtype, w)
}

func (cfg *Writers) Get(typ reflect2.Type) (obj Writer, err error) {
	rtype := uint64(typ.RType())
	if cached, exit := cfg.cache.Load(rtype); exit {
		obj = cached.(Writer)
		return
	}
	if cached, exit := cfg.processing.Load(rtype); exit {
		obj = cached.(Writer)
		return
	}
	var groupKeyBytes []byte
	cachedGroupKey := cfg.groupKeyPool.Get()
	if cachedGroupKey != nil {
		groupKeyBytes = cachedGroupKey.([]byte)
	} else {
		groupKeyBytes = make([]byte, 8)
	}
	binary.LittleEndian.PutUint64(groupKeyBytes, rtype)
	groupKey := unsafe.String(unsafe.SliceData(groupKeyBytes), len(groupKeyBytes))
	v, vErr, _ := cfg.group.Do(groupKey, func() (interface{}, error) {
		vv, objErr := WriterOf(cfg, typ)
		if objErr != nil {
			return nil, objErr
		}
		cfg.cache.Store(rtype, vv)
		return vv, nil
	})
	if vErr != nil {
		err = vErr
		return
	}
	cfg.group.Forget(groupKey)
	cfg.groupKeyPool.Put(groupKeyBytes)
	obj = v.(Writer)
	return
}

func (cfg *Writers) addProcessing(typ reflect2.Type, w Writer) {
	rtype := uint64(typ.RType())
	cfg.processing.Store(rtype, w)
}

func (cfg *Writers) removeProcessing(typ reflect2.Type) {
	rtype := uint64(typ.RType())
	cfg.processing.Delete(rtype)
}

func WriterOf(cfg *Writers, typ reflect2.Type) (v Writer, err error) {
	switch typ.Kind() {
	case reflect.String:
		v = NewStringWriter()
		break
	case reflect.Bool:
		v = NewBoolWriter()
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = NewIntWriter()
		break
	case reflect.Float32, reflect.Float64:
		v = NewFloatWriter()
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v = NewUintWriter()
		break
	case reflect.Struct:
		v, err = NewStruct(cfg, typ)
		break
	case reflect.Ptr:
		typ = typ.(reflect2.PtrType).Elem()
		v, err = WriterOf(cfg, typ)
		break
	case reflect.Slice:
		v, err = NewSliceType(cfg, typ.(reflect2.SliceType))
		break
	case reflect.Map:
		v, err = NewMapType(cfg, typ.(reflect2.MapType))
		break
	default:
		err = fmt.Errorf("copier: not support %s dst type", typ.String())
		return
	}
	return
}
