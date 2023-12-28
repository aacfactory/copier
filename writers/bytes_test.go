package writers_test

import (
	"bytes"
	"github.com/aacfactory/copier/writers"
	"github.com/modern-go/reflect2"
	"testing"
)

func TestBytesWriter_Write(t *testing.T) {
	writer := writers.NewBytesWriter(reflect2.TypeOf([]byte{}))
	dst := make([]byte, 0)
	src := []byte("0123456789")
	err := writer.Write(reflect2.PtrOf(&dst), reflect2.PtrOf(&src), reflect2.TypeOf(src))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(bytes.Equal(dst, src))
}

func TestBytesWriter_String(t *testing.T) {
	writer := writers.NewBytesWriter(reflect2.TypeOf([]byte{}))
	dst := make([]byte, 0)
	src := "0123456789"
	err := writer.Write(reflect2.PtrOf(&dst), reflect2.PtrOf(&src), reflect2.TypeOf(src))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(bytes.Equal(dst, []byte(src)))
}

type BytesText struct {
	p []byte
}

func (b *BytesText) UnmarshalText(text []byte) error {
	b.p = text
	return nil
}

func (b BytesText) MarshalText() (text []byte, err error) {
	return b.p, err
}

func TestBytesWriter_Text(t *testing.T) {
	writer := writers.NewBytesWriter(reflect2.TypeOf([]byte{}))
	dst := make([]byte, 0)
	src := BytesText{
		p: []byte("0123456789"),
	}
	err := writer.Write(reflect2.PtrOf(&dst), reflect2.PtrOf(&src), reflect2.TypeOf(src))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(bytes.Equal(dst, src.p))
}
