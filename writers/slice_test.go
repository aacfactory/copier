package writers_test

import (
	"github.com/aacfactory/copier/writers"
	"github.com/modern-go/reflect2"
	"testing"
)

func TestSliceWriter_Write(t *testing.T) {
	dst := make([]int, 0)
	writer, _ := writers.NewSliceType(writers.New("copy"), reflect2.TypeOf(dst).(reflect2.SliceType))
	src := []int{0, 1, 2}
	err := writer.Write(&dst, &src)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst)
}

func TestSliceWriter_WriteString(t *testing.T) {
	dst := make([]int, 0)
	writer, _ := writers.NewSliceType(writers.New("copy"), reflect2.TypeOf(dst).(reflect2.SliceType))
	src := []string{"0", "1", "2"}
	err := writer.Write(&dst, &src)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst)
}
