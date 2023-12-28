package copier

import (
	"errors"
	"fmt"
	"github.com/aacfactory/copier/writers"
	"github.com/modern-go/reflect2"
	"reflect"
)

var (
	DefaultConfig = Config{
		TagKey: "copy",
	}.Freeze()
)

type Config struct {
	TagKey  string
	Writers []writers.Writer
}

func (c Config) Register(w writers.Writer) {
	c.Writers = append(c.Writers, w)
}

func (c Config) Freeze() API {
	api := &frozenConfig{
		writers: writers.New(c.TagKey, c.Writers...),
	}
	return api
}

type API interface {
	Copy(dst any, src any) (err error)
}

type frozenConfig struct {
	writers *writers.Writers
}

func (cfg *frozenConfig) Copy(dst any, src any) (err error) {

	if src == nil {
		err = fmt.Errorf("copier: src must not be nil")
		return
	}
	if dst == nil {
		err = fmt.Errorf("copier: dst must not be nil")
		return
	}

	dstType := reflect2.TypeOf(dst)
	if dstType.Kind() != reflect.Ptr {
		err = fmt.Errorf("copier: dst must be ptr")
		return
	}
	dstPtr := reflect2.PtrOf(dst)

	srcType := reflect2.TypeOf(src)
	srcPtr := reflect2.PtrOf(src)

	if dstType.RType() == srcType.RType() {
		dstType.UnsafeSet(dstPtr, srcPtr)
		return
	}

	dstObj, dstErr := cfg.writers.Get(dstType)
	if dstErr != nil {
		err = errors.Join(errors.New("copier: copy failed"), dstErr)
		return
	}
	if err = dstObj.Write(dstPtr, srcPtr, srcType); err != nil {
		err = errors.Join(errors.New("copier: copy failed"), err)
		return
	}
	return
}
