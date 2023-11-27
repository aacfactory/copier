package copier_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aacfactory/copier"
	"reflect"
	"testing"
	"time"
)

type Date time.Time

type User struct {
	Name string
}

type NullJson[E any] struct {
	Value E
	Valid bool
}

func (n NullJson[E]) Scan(src any) error {
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
		SQLTime:   sql.NullTime{Time: time.Now()},
		SQLString: sql.NullString{String: "x"},
	}
	err := copier.Copy(foo, bar)
	fmt.Println(err)
	fmt.Println(fmt.Sprintf("%+v", foo))
	fmt.Println(foo.Bazz)
	fmt.Println(foo.Bazs[0])
	fmt.Println(foo.MM["a"])
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
		"a": &Baz{
			X: "b",
		},
	}
	err := copier.Copy(&dst, src)
	fmt.Println(err)
	fmt.Println(dst)
	fmt.Println(dst["a"])
}

func TestValueOf(t *testing.T) {
	bar := Bar{
		User: NullJson[User]{
			Value: User{Name: "username"},
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
		SQLString: sql.NullString{String: "x", Valid: false},
	}
	dst, err := copier.ValueOf[Foo](bar)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(fmt.Sprintf("%+v", dst))
}

type XUser struct {
	NullJson[User]
}

func TestNullUser(t *testing.T) {
	nu := XUser{}
	rv := reflect.ValueOf(nu)
	for i := 0; i < rv.NumField(); i++ {
		fmt.Println(rv.Type().Field(i))
	}
}

type Internal struct {
	Id   string
	Name string
	s    string
}

type SI struct {
	*Internal
	Bar string
}

type SS struct {
	Id   string
	Name string
	Bar  string
	s    string
}

func TestValueOf2(t *testing.T) {
	si, siErr := copier.ValueOf[SI](SS{
		Id:   "1",
		Name: "name",
		Bar:  "bar",
		s:    "s",
	})
	if siErr != nil {
		fmt.Println(siErr)
		return
	}
	fmt.Println(fmt.Sprintf("%+v", si))

	ss, ssErr := copier.ValueOf[SS](si)
	if ssErr != nil {
		fmt.Println(ssErr)
		return
	}
	fmt.Println(fmt.Sprintf("%+v", ss))

	si1, si1Err := copier.ValueOf[SI](si)
	if si1Err != nil {
		fmt.Println(si1Err)
		return
	}
	fmt.Println(fmt.Sprintf("%+v", si1))
}
