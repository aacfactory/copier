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
		TagKey: "copier",
	}.Freeze()
)

type Config struct {
	TagKey string
}

func (c Config) Freeze() API {
	api := &frozenConfig{
		config:  c,
		writers: writers.New(c.TagKey),
	}
	return api
}

// API represents a frozen Config.
type API interface {
	Copy(dst any, src any) (err error)
}

type frozenConfig struct {
	config  Config
	writers *writers.Writers
}

func (cfg *frozenConfig) Copy(dst any, src any) (err error) {
	dstType := reflect2.TypeOf(dst)
	if dstType.Kind() != reflect.Ptr {
		err = fmt.Errorf("copier: dst must be ptr")
		return
	}
	if reflect2.IsNil(dst) {
		err = fmt.Errorf("copier: dst must not be nil")
		return
	}
	if src == nil {
		err = fmt.Errorf("copier: src must not be nil")
		return
	}

	dstObj, dstErr := cfg.writers.Get(dstType)
	if dstErr != nil {
		err = errors.Join(errors.New("copier: copy failed"), dstErr)
		return
	}
	if err = dstObj.Write(reflect2.PtrOf(dst), reflect2.PtrOf(src), reflect2.TypeOf(src)); err != nil {
		err = errors.Join(errors.New("copier: copy failed"), err)
		return
	}
	return
}
