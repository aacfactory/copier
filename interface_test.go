package copier_test

import (
	"github.com/aacfactory/copier"
	"testing"
	"time"
)

type Entry interface {
	String() string
}

type EntryImpl struct {
	s string
}

func (e *EntryImpl) String() string {
	return e.s
}

func TestCopy_Interface(t *testing.T) {
	var dst Entry
	var src = EntryImpl{s: time.Now().String()}
	err := copier.Copy(&dst, src)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst.String())
}

func TestCopy_ToInterface(t *testing.T) {
	var dst = EntryImpl{}
	var src Entry
	src = &EntryImpl{s: time.Now().String()}
	err := copier.Copy(&dst, src)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst.String())
}
