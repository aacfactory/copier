package copier_test

import (
	"github.com/aacfactory/copier"
	"testing"
	"time"
)

func TestCopy_Slice(t *testing.T) {
	foo := make([]Foo, 0, 1)
	bar := []Bar{
		{
			String: "1",
			Bool:   true,
			Int:    1,
			Float:  2.2,
			Uint:   3,
			Bytes:  []byte("bytes"),
			Time:   time.Now(),
			Slice:  []BarElem{{String: "1", Int64: 1}, {String: "2", Int64: 2}},
			Map:    map[int]BarElem{1: {String: "1", Int64: 1}, 2: {String: "2", Int64: 2}},
		},
		{
			String: "2",
			Bool:   true,
			Int:    1,
			Float:  2.2,
			Uint:   3,
			Bytes:  []byte("bytes"),
			Time:   time.Now(),
			Slice:  []BarElem{{String: "1", Int64: 1}, {String: "2", Int64: 2}},
			Map:    map[int]BarElem{1: {String: "1", Int64: 1}, 2: {String: "2", Int64: 2}},
		},
	}
	err := copier.Copy(&foo, bar)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", foo)
}

func TestValueOf_Slice(t *testing.T) {
	bar := []Bar{
		{
			String: "1",
			Bool:   true,
			Int:    1,
			Float:  2.2,
			Uint:   3,
			Bytes:  []byte("bytes"),
			Time:   time.Now(),
			Slice:  []BarElem{{String: "1", Int64: 1}, {String: "2", Int64: 2}},
			Map:    map[int]BarElem{1: {String: "1", Int64: 1}, 2: {String: "2", Int64: 2}},
		},
		{
			String: "2",
			Bool:   true,
			Int:    1,
			Float:  2.2,
			Uint:   3,
			Bytes:  []byte("bytes"),
			Time:   time.Now(),
			Slice:  []BarElem{{String: "1", Int64: 1}, {String: "2", Int64: 2}},
			Map:    map[int]BarElem{1: {String: "1", Int64: 1}, 2: {String: "2", Int64: 2}},
		},
	}
	foo, err := copier.ValueOf[[]Foo](bar)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", foo)
}
