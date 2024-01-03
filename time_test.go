package copier_test

import (
	"github.com/aacfactory/copier"
	"testing"
	"time"
)

type IntTime struct {
	Int time.Time
}

func TestCopy_time(t *testing.T) {
	foo := Foo{}
	bar := IntTime{
		Int: time.Now(),
	}
	err := copier.Copy(&foo, bar)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", foo.Int)
}

func TestCopy_time_from_ine(t *testing.T) {
	foo := Foo{
		Int: 1704272255419,
	}
	bar := IntTime{}
	err := copier.Copy(&bar, foo)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", bar.Int)
}
