package copier_test

import (
	"github.com/aacfactory/copier"
	"testing"
	"time"
)

type BarBytes struct {
	String []byte
}

func TestCopy_string_bytes(t *testing.T) {
	foo := Foo{}
	bar := BarBytes{
		String: []byte("string"),
	}
	err := copier.Copy(&foo, bar)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", foo)
}

type BarTime struct {
	String time.Time
}

func TestCopy_string_time(t *testing.T) {
	foo := Foo{}
	bar := BarTime{
		String: time.Now(),
	}
	err := copier.Copy(&foo, bar)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", foo)
}

type BarText struct {
	String Text
}

func TestCopy_string_text(t *testing.T) {
	foo := Foo{}
	bar := BarText{
		String: Text{
			p: []byte("string"),
		},
	}
	err := copier.Copy(&foo, bar)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", foo)
}
