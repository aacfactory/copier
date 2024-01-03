package copier_test

import (
	"database/sql"
	"github.com/aacfactory/copier"
	"testing"
	"time"
)

type BarSQL struct {
	String sql.NullString
	Bool   sql.NullBool
	Int    sql.NullInt64
	Float  sql.NullFloat64
	Time   sql.NullTime
	Byte   sql.NullByte
}

func TestCopy_Sql(t *testing.T) {
	foo := Foo{}
	bar := BarSQL{
		String: sql.NullString{
			String: "string",
			Valid:  true,
		},
		Bool: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
		Int: sql.NullInt64{
			Int64: 1,
			Valid: true,
		},
		Float: sql.NullFloat64{
			Float64: 2.2,
			Valid:   true,
		},
		Time: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Byte: sql.NullByte{
			Byte:  1,
			Valid: true,
		},
	}
	err := copier.Copy(&foo, bar)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", foo)
}

func TestCopy_Sql2(t *testing.T) {
	foo := Foo{
		String: "",
		Bool:   true,
		Int:    1,
		Float:  2.2,
		Uint:   0,
		Byte:   'F',
		Bytes:  nil,
		Time:   time.Now(),
		Slice:  nil,
		Map:    nil,
	}
	bar := BarSQL{}
	err := copier.Copy(&bar, foo)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", bar)
}
