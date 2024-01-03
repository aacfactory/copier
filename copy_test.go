package copier_test

import (
	"github.com/aacfactory/copier"
	"testing"
	"time"
)

type FooElem struct {
	String string
	Int    int `copy:"int"`
}

type Foo struct {
	String string
	Bool   bool
	Int    int
	Float  float64
	Uint   uint
	Byte   byte
	Bytes  []byte
	Time   time.Time
	Slice  []FooElem
	Map    map[string]FooElem
}

type BarElem struct {
	String string
	Int64  int `copy:"int"`
}

type Bar struct {
	String string
	Bool   bool
	Int    int
	Float  float64
	Uint   uint
	Byte   byte
	Bytes  []byte
	Time   time.Time
	Slice  []BarElem
	Map    map[int]BarElem
}

func TestCopy(t *testing.T) {
	foo := Foo{}
	bar := Bar{
		String: "string",
		Bool:   true,
		Int:    1,
		Float:  2.2,
		Uint:   3,
		Byte:   byte('B'),
		Bytes:  []byte("bytes"),
		Time:   time.Now(),
		Slice:  []BarElem{{String: "1", Int64: 1}, {String: "2", Int64: 2}},
		Map:    map[int]BarElem{1: {String: "1", Int64: 1}, 2: {String: "2", Int64: 2}},
	}
	err := copier.Copy(&foo, bar)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", foo)
}

func TestValueOf(t *testing.T) {
	bar := Bar{
		String: "string",
		Bool:   true,
		Int:    1,
		Float:  2.2,
		Uint:   3,
		Byte:   byte('B'),
		Bytes:  []byte("bytes"),
		Time:   time.Now(),
		Slice:  []BarElem{{String: "1", Int64: 1}, {String: "2", Int64: 2}},
		Map:    map[int]BarElem{1: {String: "1", Int64: 1}, 2: {String: "2", Int64: 2}},
	}
	foo, err := copier.ValueOf[Foo](bar)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", foo)
}

func BenchmarkCopy(b *testing.B) {
	b.ReportAllocs()
	bar := Bar{
		String: "string",
		Bool:   true,
		Int:    1,
		Float:  2.2,
		Uint:   3,
		Bytes:  []byte("bytes"),
		Time:   time.Now(),
		Slice:  []BarElem{{String: "1", Int64: 1}, {String: "2", Int64: 2}},
		Map:    map[int]BarElem{1: {String: "1", Int64: 1}, 2: {String: "2", Int64: 2}},
	}
	for i := 0; i < b.N; i++ {
		foo := Foo{}
		err := copier.Copy(&foo, bar)
		if err != nil {
			b.Error(err)
			return
		}
	}
}
