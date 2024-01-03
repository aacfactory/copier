package copier_test

import (
	"github.com/aacfactory/copier"
	"testing"
)

type Text struct {
	p []byte
}

func (t *Text) UnmarshalText(text []byte) error {
	t.p = text
	return nil
}

func (t Text) MarshalText() (text []byte, err error) {
	text = t.p
	return
}

type FooText struct {
	String Text
}

func TestCopy_text(t *testing.T) {
	foo := FooText{}
	bar := Bar{
		String: "string",
	}
	err := copier.Copy(&foo, bar)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", foo.String)
}
