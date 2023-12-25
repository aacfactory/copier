package descriptors_test

import (
	"database/sql"
	"github.com/aacfactory/copier/descriptors"
	"github.com/modern-go/reflect2"
	"testing"
)

type NullString struct {
	sql.NullString
}

func TestNullStringDescriptor_ValueOf(t *testing.T) {
	desc, descErr := descriptors.DescribeNullString(reflect2.TypeOf(new(sql.NullString)))
	if descErr != nil {
		t.Error(descErr)
		return
	}
	v := NullString{
		sql.NullString{
			String: "ss",
			Valid:  true,
		},
	}
	ptr, valid := desc.ValueOf(reflect2.PtrOf(v))
	if valid {
		s := *(*string)(ptr)
		t.Log(s)
	}
}
