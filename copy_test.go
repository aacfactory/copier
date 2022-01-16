package copier_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aacfactory/copier"
	"testing"
	"time"
)

type Date time.Time

type Foo struct {
	Str       string
	Time      Date `copy:"Time"`
	I         int  `copy:"i"`
	Bytes     json.RawMessage
	Baz       Baz
	Bazs      []Faz
	IS        []int
	ISS       [][]int
	MM        map[string]*Faz
	SQLTime   time.Time
	SQLString string
}

type Bar struct {
	Str       string
	Now       time.Time `copy:"Time"`
	X         int64     `copy:"i"`
	Bytes     []byte
	Baz       Baz
	Bazs      []*Baz
	IS        []int
	ISS       [][]int
	MM        map[string]*Baz
	SQLTime   sql.NullTime
	SQLString sql.NullString
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
		MM: map[string]*Baz{
			"a": {X: "a"},
			"b": {X: "b"},
			"c": {X: "c"},
		},
		SQLTime:   sql.NullTime{Time: time.Now()},
		SQLString: sql.NullString{String: "x"},
	}
	err := copier.Copy(foo, bar)
	fmt.Println(err)
	fmt.Println(fmt.Sprintf("%+v", foo))

}
