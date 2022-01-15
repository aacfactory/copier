package copier_test

import (
	"fmt"
	"github.com/aacfactory/copier"
	"testing"
	"time"
)

type Date time.Time

type Foo struct {
	Str  string
	Time Date `copy:"Time"`
	i    int
}

type Bar struct {
	Str string
	Now time.Time `copy:"Time"`
}

func TestCopy(t *testing.T) {
	foo := &Foo{}
	bar := Bar{
		Str: "str",
		Now: time.Now(),
	}
	err := copier.Copy(foo, bar)
	fmt.Println(err)
	fmt.Println(foo)
}
