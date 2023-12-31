package writers_test

import (
	"database/sql"
	"github.com/aacfactory/copier/writers"
	"testing"
)

func TestBoolWriter_Write(t *testing.T) {
	writer := writers.NewBoolWriter()
	dst := false
	src := true
	err := writer.Write(&dst, &src)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst)
}

func TestBoolWriter_WriteSQL(t *testing.T) {
	writer := writers.NewBoolWriter()
	dst := false
	src := sql.NullBool{
		Bool:  true,
		Valid: true,
	}
	err := writer.Write(&dst, &src)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst)
}

func TestBoolWriter_WriteSQLPtr(t *testing.T) {
	writer := writers.NewBoolWriter()
	dst := false
	src := &sql.NullBool{
		Bool:  true,
		Valid: true,
	}
	err := writer.Write(&dst, &src)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst)
}

func TestBoolWriter_WriteString(t *testing.T) {
	writer := writers.NewBoolWriter()
	dst := false
	src := "true"
	err := writer.Write(&dst, src)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst)
}

func TestBoolWriter_WriteInt(t *testing.T) {
	writer := writers.NewBoolWriter()
	dst := false
	src := 1
	err := writer.Write(&dst, &src)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst)
}

func TestBoolWriter_WriteByte(t *testing.T) {
	writer := writers.NewBoolWriter()
	dst := false
	src := 't'
	err := writer.Write(&dst, &src)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst)
}
