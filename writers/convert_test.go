package writers_test

import (
	"github.com/aacfactory/copier/writers"
	"github.com/modern-go/reflect2"
	"testing"
	"time"
)

func (w *Date) Convert() any {
	return time.Date(w.Year, w.Month, w.Day, 0, 0, 0, 0, time.Local)
}

func TestIsConvertible(t *testing.T) {
	w := writers.NewTimeWriter(reflect2.TypeOf(time.Time{}))
	dst := time.Time{}
	src := Date{
		Year:  time.Now().Year(),
		Month: time.Now().Month(),
		Day:   time.Now().Day(),
	}
	err := w.Write(&dst, &src)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst)
	t.Log(src.Year == dst.Year(), src.Month == dst.Month(), src.Day == dst.Day())
}
