package copier_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/aacfactory/copier"
	"testing"
	"time"
)

type Date time.Time

type User struct {
	Name string
}

type NullJson[E any] struct {
	E     E
	Valid bool
}

func (n NullJson[E]) Value() (driver.Value, error) {
	return n.E, nil
}

func (n *NullJson[E]) Scan(src any) error {
	return nil
}

type Foo struct {
	User      User
	Str       string
	IS        []int
	ISS       [][]int
	Time      Date `copy:"Time"`
	I         int  `copy:"i"`
	SQLTime   time.Time
	SQLString string
	Bytes     json.RawMessage
	Baz       Faz
	Bazz      *Faz
	Bazs      []Faz
	Bazzs     []*Faz
	MM        map[string]*Baz
}

type Bar struct {
	User      NullJson[User]
	Str       string
	Now       time.Time `copy:"Time"`
	X         int64     `copy:"i"`
	Bytes     []byte
	Baz       Baz
	Bazz      *Baz
	Bazs      []Baz
	Bazzs     []*Baz
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
		User: NullJson[User]{
			E: User{
				Name: "name",
			},
			Valid: true,
		},
		Str:   "str",
		Now:   time.Now(),
		X:     100,
		Bytes: []byte(`{"a":1}`),
		Baz:   Baz{X: "0"},
		Bazz:  &Baz{X: "1"},
		Bazs:  []Baz{{X: "1"}},
		Bazzs: []*Baz{{X: "1"}},
		IS:    []int{1, 2},
		ISS:   [][]int{{1, 2}, {3, 4}},
		MM: map[string]*Baz{
			"a": {X: "a"},
			"b": {X: "b"},
			"c": {X: "c"},
		},
		SQLTime:   sql.NullTime{Time: time.Now(), Valid: true},
		SQLString: sql.NullString{String: "x"},
	}
	err := copier.Copy(foo, bar)
	fmt.Println(err)
	fmt.Println(fmt.Sprintf("%+v", foo))
	fmt.Println(time.Time(foo.Time).String())
	fmt.Println(foo.Baz)
	fmt.Println(foo.Bazz)
	fmt.Println(len(foo.Bazs), foo.Bazs == nil, foo.Bazs)
	fmt.Println(len(foo.Bazzs), foo.Bazzs == nil, foo.Bazzs)
}

func TestArray(t *testing.T) {
	dst := make([]*Faz, 0, 1)
	src := append(make([]*Faz, 0, 1), &Faz{
		X: "d",
	})
	err := copier.Copy(&dst, src)
	fmt.Println(err)
	fmt.Println(dst)
	fmt.Println(dst[0])

}

func TestMap(t *testing.T) {
	dst := make(map[string]*Faz)
	src := map[string]*Baz{
		"a": {
			X: "b",
		},
	}
	err := copier.Copy(&dst, src)
	fmt.Println(err)
	fmt.Println(dst)
	fmt.Println(dst["a"])
}
