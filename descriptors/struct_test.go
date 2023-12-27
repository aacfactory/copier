package descriptors_test

import (
	"github.com/aacfactory/copier/descriptors"
	"github.com/modern-go/reflect2"
	"testing"
	"time"
)

type Anonymous struct {
	AnonymousKey   string
	AnonymousValue string
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
	sss     string
}

func TestStructFields_Get(t *testing.T) {
	foo := Foo{
		Anonymous: Anonymous{
			AnonymousValue: "AnonymousValue",
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
		t.Log("name:", f.Name, "type:", fType, "value:", fv, fType.UnsafeIsNil(fp))
	}
	f := desc.FieldByTag("copier", "str")
	t.Log(f)
	v := Foo{}
	vp := reflect2.PtrOf(&v)
	anoF := desc.Field("AnonymousValue")
	_, anoFT, _ := anoF.ValueOf(vp)
	anoFT.UnsafeSet(vp, reflect2.PtrOf("ano"))
	t.Log(v.Anonymous.AnonymousValue)
	t.Log("-----------")
	for _, fieldDescriptor := range desc.Fields() {
		t.Log(fieldDescriptor.Name, "---")
		for _, field := range fieldDescriptor.Field {
			t.Log(field.StructField.Name)
		}
		t.Log("***")
	}
}
