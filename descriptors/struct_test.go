package descriptors_test

import (
	"github.com/aacfactory/copier/descriptors"
	"github.com/modern-go/reflect2"
	"testing"
	"time"
)

type Anonymous struct {
	Anonymous string
}

type Bar struct {
	String string
	Next   *Bar
}

type Foo struct {
	Anonymous
	String  string `copier:"str"`
	Boolean bool
	Int     int
	Long    int64
	Float   float32
	Double  float64
	Uint    uint64
	Time    time.Time
	Dur     time.Duration
	Byte    byte
	Bytes   []byte
	Bar     Bar
	Baz     *Bar
	Bars    []Bar
	Map     map[string]Bar
}

func TestStructFields_Get(t *testing.T) {
	foo := Foo{
		Anonymous: Anonymous{
			Anonymous: "Anonymous",
		},
		String:  "xxx",
		Boolean: false,
		Int:     0,
		Long:    0,
		Float:   0,
		Double:  0,
		Uint:    0,
		Time:    time.Time{},
		Dur:     0,
		Byte:    0,
		Bytes:   nil,
		Bar:     Bar{},
		Baz:     nil,
		Bars:    nil,
		Map:     nil,
	}
	ptr := reflect2.PtrOf(&foo)
	desc := descriptors.DescribeStruct(reflect2.TypeOf(Foo{}))
	for _, f := range desc.Fields() {
		fp, fType, fErr := f.ValueOf(ptr)
		if fErr != nil {
			t.Error("field:", f.Name, fErr)
			return
		}
		fv := fType.UnsafeIndirect(fp)
		t.Log("name:", f.Name, "type:", fType, "value:", fv)
	}
	f := desc.FieldByTag("copier", "str")
	t.Log(f)
}
