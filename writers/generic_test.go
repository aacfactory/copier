package writers_test

import (
	"database/sql"
	"fmt"
	"github.com/aacfactory/copier/writers"
	"github.com/modern-go/reflect2"
	"testing"
	"time"
	"unsafe"
)

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

func (w *Date) UnsafeWrite(srcPtr unsafe.Pointer, srcType reflect2.Type) (err error) {
	if srcType.RType() == reflect2.RTypeOf(time.Time{}) || srcType.Type1().ConvertibleTo(reflect2.TypeOf(time.Time{}).Type1()) {
		src := time.Time{}
		reflect2.TypeOf(src).UnsafeSet(reflect2.PtrOf(&src), srcPtr)
		w.Year = src.Year()
		w.Month = src.Month()
		w.Day = src.Day()
		return
	}
	if srcType.RType() == reflect2.RTypeOf(sql.NullTime{}) || srcType.Type1().ConvertibleTo(reflect2.TypeOf(sql.NullTime{}).Type1()) {
		src := sql.NullTime{}
		reflect2.TypeOf(src).UnsafeSet(reflect2.PtrOf(&src), srcPtr)
		if src.Valid {
			w.Year = src.Time.Year()
			w.Month = src.Time.Month()
			w.Day = src.Time.Day()
		}
		return
	}
	err = fmt.Errorf("unsupported %s", srcType.String())
	return
}

func TestGenericInterfaceWriter_Write(t *testing.T) {
	w := writers.NewGenericWriter(reflect2.TypeOf(Date{}))
	dst := Date{}
	src := time.Now()
	err := w.Write(reflect2.PtrOf(&dst), reflect2.PtrOf(&src), reflect2.TypeOf(src))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst.Year == src.Year(), dst.Month == src.Month(), dst.Day == src.Day())
}

func TestGenericInterfaceWriter_WriteSQL(t *testing.T) {
	w := writers.NewGenericWriter(reflect2.TypeOf(Date{}))
	dst := Date{}
	src := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	err := w.Write(reflect2.PtrOf(&dst), reflect2.PtrOf(&src), reflect2.TypeOf(src))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst.Year == src.Time.Year(), dst.Month == src.Time.Month(), dst.Day == src.Time.Day())
}
