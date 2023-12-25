package writers

import (
	"encoding/binary"
	"fmt"
	"github.com/modern-go/reflect2"
	"golang.org/x/sync/singleflight"
	"reflect"
	"sync"
	"unsafe"
)

type Writer interface {
	Write(dstPtr unsafe.Pointer, reader Reader) (err error)
}

func New(tagKey string) *Writers {
	w := &Writers{
		values:     sync.Map{},
		processing: sync.Map{},
		groupKeys:  sync.Pool{},
		group:      singleflight.Group{},
		tagKey:     tagKey,
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
	return w
}

type Writers struct {
	values     sync.Map
	processing sync.Map
	groupKeys  sync.Pool
	group      singleflight.Group
	tagKey     string
}

func (writers *Writers) Get(typ reflect2.Type) (w Writer, err error) {
	rtype := typ.RType()
	cached, has := writers.load(rtype)
	if has {
		w = cached
		return
	}
	var groupKey []byte
	cachedGroupKey := writers.groupKeys.Get()
	if cachedGroupKey == nil {
		groupKey = make([]byte, 8)
	} else {
		groupKey = cachedGroupKey.([]byte)
	}
	binary.LittleEndian.PutUint64(groupKey, uint64(rtype))
	groupKeyStr := unsafe.String(unsafe.SliceData(groupKey), 8)
	v, doErr, _ := writers.group.Do(groupKeyStr, func() (v any, err error) {
		nw, nErr := writers.create(typ, rtype)
		if nErr != nil {
			err = nErr
			return
		}
		writers.set(rtype, nw)
		v = nw
		return
	})
	writers.group.Forget(groupKeyStr)
	writers.groupKeys.Put(groupKey)
	if doErr != nil {
		err = doErr
		return
	}
	w = v.(Writer)
	return
}

func (writers *Writers) Register(typ reflect2.Type, w Writer) {
	rtype := typ.RType()
	writers.set(rtype, w)
}

func (writers *Writers) load(rtype uintptr) (w Writer, has bool) {
	v, exist := writers.values.Load(uint64(rtype))
	if exist {
		w, has = v.(Writer)
		return
	}
	return
}

func (writers *Writers) set(rtype uintptr, w Writer) {
	writers.values.Store(uint64(rtype), w)
}

func (writers *Writers) loadProcessing(rtype uintptr) (w Writer, has bool) {
	v, exist := writers.processing.Load(uint64(rtype))
	if exist {
		w, has = v.(Writer)
		return
	}
	return
}

func (writers *Writers) setProcessing(rtype uintptr, w Writer) {
	writers.processing.Store(uint64(rtype), w)
}

func (writers *Writers) create(typ reflect2.Type, rtype uintptr) (w Writer, err error) {
	if processing, inProcessing := writers.loadProcessing(rtype); inProcessing {
		w = processing
		return
	}
	if created, inCreated := writers.load(rtype); inCreated {
		w = created
		return
	}
	switch typ.Kind() {
	case reflect.String:
		w = NewStringWriter()
		break
	case reflect.Bool:
		w = NewBoolWriter()
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		w = NewIntWriter()
		break
	case reflect.Float32, reflect.Float64:
		w = NewFloatWriter()
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		w = NewUintWriter()
		break
	case reflect.Struct:

		break
	case reflect.Slice:

		break
	case reflect.Map:

		break
	default:
		err = fmt.Errorf("copy failed for %v is not supported", typ.Kind())
		return
	}
	return
}
