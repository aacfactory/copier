package copier_test

import (
	"github.com/aacfactory/copier"
	"testing"
)

type Anonymous struct {
	AnoKey string
	AnoVal string
}

type FooAno struct {
	*Anonymous
	String string
}

type BarAno struct {
	Anonymous
	String string
}

func TestCopy_Ano(t *testing.T) {
	foo := FooAno{}
	bar := BarAno{
		String: "string",
		Anonymous: Anonymous{
			AnoKey: "ano_key",
			AnoVal: "ano_val",
		},
	}
	err := copier.Copy(&foo, bar)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", foo)
	t.Logf("%+v", foo.Anonymous)
}
