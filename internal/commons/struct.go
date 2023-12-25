package commons

import (
	"github.com/modern-go/reflect2"
	"reflect"
)

type StructField struct {
	Name  string
	Field []*reflect2.UnsafeStructField
	anon  *reflect2.UnsafeStructType
}

type StructFields []*StructField

type StructDescriptor struct {
	Type   reflect2.Type
	Fields StructFields
}

func DescribeStruct(tagKey string, typ reflect2.Type) *StructDescriptor {
	structType := typ.(*reflect2.UnsafeStructType)
	fields := StructFields{}

	var curr []StructField
	next := []StructField{{anon: structType}}

	visited := map[uintptr]bool{}

	for len(next) > 0 {
		curr, next = next, curr[:0]

		for _, f := range curr {
			rtype := f.anon.RType()
			if visited[f.anon.RType()] {
				continue
			}
			visited[rtype] = true

			for i := 0; i < f.anon.NumField(); i++ {
				field := f.anon.Field(i).(*reflect2.UnsafeStructField)
				isUnexported := field.PkgPath() != ""

				chain := make([]*reflect2.UnsafeStructField, len(f.Field)+1)
				copy(chain, f.Field)
				chain[len(f.Field)] = field

				if field.Anonymous() {
					t := field.Type()
					if t.Kind() == reflect.Ptr {
						t = t.(*reflect2.UnsafePtrType).Elem()
					}
					if t.Kind() != reflect.Struct {
						continue
					}

					next = append(next, StructField{Field: chain, anon: t.(*reflect2.UnsafeStructType)})
					continue
				}

				// Ignore unexported fields.
				if isUnexported {
					continue
				}

				fieldName := field.Name()
				if tagKey != "" {
					if tag, ok := field.Tag().Lookup(tagKey); ok {
						fieldName = tag
					}
				}

				fields = append(fields, &StructField{
					Name:  fieldName,
					Field: chain,
				})
			}
		}
	}

	return &StructDescriptor{
		Type:   structType,
		Fields: fields,
	}
}
