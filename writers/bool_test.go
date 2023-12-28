package writers_test

import (
	"database/sql"
	"github.com/aacfactory/copier/writers"
	"github.com/modern-go/reflect2"
	"testing"
)

func TestBoolWriter_Write(t *testing.T) {
	writer := writers.NewBoolWriter()
	dst := false
	src := true
	err := writer.Write(reflect2.PtrOf(&dst), reflect2.PtrOf(&src), reflect2.TypeOf(src))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst == src)
}

func TestBoolWriter_WriteSQL(t *testing.T) {
	writer := writers.NewBoolWriter()
	dst := false
	src := sql.NullBool{
		Bool:  true,
		Valid: true,
	}
	err := writer.Write(reflect2.PtrOf(&dst), reflect2.PtrOf(&src), reflect2.TypeOf(src))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst == src.Bool)
}

func TestBoolWriter_WriteString(t *testing.T) {
	writer := writers.NewBoolWriter()
	dst := false
	src := "true"
	err := writer.Write(reflect2.PtrOf(&dst), reflect2.PtrOf(&src), reflect2.TypeOf(src))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst == true)
}

func TestBoolWriter_WriteInt(t *testing.T) {
	writer := writers.NewBoolWriter()
	dst := false
	src := 1
	err := writer.Write(reflect2.PtrOf(&dst), reflect2.PtrOf(&src), reflect2.TypeOf(src))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst == true)
}

func TestBoolWriter_WriteByte(t *testing.T) {
	writer := writers.NewBoolWriter()
	dst := false
	src := 't'
	err := writer.Write(reflect2.PtrOf(&dst), reflect2.PtrOf(&src), reflect2.TypeOf(src))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst == true)
}
