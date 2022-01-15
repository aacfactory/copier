package copier_test

import (
	"encoding/json"
	"fmt"
	"github.com/aacfactory/copier"
	"reflect"
	"testing"
	"time"
)

type Date time.Time

type Foo struct {
	Str   string
	Time  Date `copy:"Time"`
	I     int  `copy:"i"`
	Bytes json.RawMessage
	Baz   Baz
	Bazs  []Faz
	IS    []int
	ISS   [][]int
	MM    map[string]*Faz
}

type Bar struct {
	Str   string
	Now   time.Time `copy:"Time"`
	X     int64     `copy:"i"`
	Bytes []byte
	Baz   Baz
	Bazs  []*Baz
	IS    []int
	ISS   [][]int
	MM    map[string]*Baz
}

type Baz struct {
	X string
}

type Faz struct {
	X string
}

func TestCopy(t *testing.T) {
	foo := &Foo{}
	bar := Bar{
		Str:   "str",
		Now:   time.Now(),
		X:     100,
		Bytes: []byte(`{"a":1}`),
		Baz: Baz{
			X: "baz",
		},
		Bazs: []*Baz{{X: "1"}},
		IS:   []int{1, 2},
		ISS:  [][]int{{1, 2}, {3, 4}},
	}
	err := copier.Copy(foo, bar)
	fmt.Println(err)
	fmt.Println(foo)
	fmt.Println(reflect.Value{}.IsValid())
}
