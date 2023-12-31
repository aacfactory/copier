package writers_test

import (
	"github.com/aacfactory/copier/writers"
	"github.com/modern-go/reflect2"
	"testing"
)

func TestMapWriter_WriteSame(t *testing.T) {
	dst := make(map[int]int)
	src := map[int]int{1: 1, 2: 2}
	writer, _ := writers.NewMapWriter(writers.New("copy"), reflect2.TypeOf(dst).(reflect2.MapType))

	err := writer.Write(&dst, &src)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst)
}

func TestMapWriter_Write(t *testing.T) {
	dst := make(map[string]string)
	src := map[int]int{1: 1, 2: 2}
	writer, _ := writers.NewMapWriter(writers.New("copy"), reflect2.TypeOf(dst).(reflect2.MapType))

	err := writer.Write(&dst, &src)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst)
}
