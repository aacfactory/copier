package copier

import (
	"github.com/modern-go/reflect2"
	"reflect"
)

type Descriptor interface {
	Type() reflect2.Type
}

type StructField struct {
	name  string
	field []*reflect2.UnsafeStructField
	anon  *reflect2.UnsafeStructType
}

type StructFields []*StructField

func describeStruct(tagKey string, typ reflect2.Type) Descriptor {
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

				chain := make([]*reflect2.UnsafeStructField, len(f.field)+1)
				copy(chain, f.field)
				chain[len(f.field)] = field

				if field.Anonymous() {
					t := field.Type()
					if t.Kind() == reflect.Ptr {
						t = t.(*reflect2.UnsafePtrType).Elem()
					}
					if t.Kind() != reflect.Struct {
						continue
					}

					next = append(next, StructField{field: chain, anon: t.(*reflect2.UnsafeStructType)})
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
					name:  fieldName,
					field: chain,
				})
			}
		}
	}

	return &StructDescriptor{
		typ:    structType,
		fields: fields,
	}
}

type StructDescriptor struct {
	typ    reflect2.Type
	fields StructFields
}

func (desc *StructDescriptor) Type() reflect2.Type {
	return desc.typ
}

func (desc *StructDescriptor) Fields() StructFields {
	return desc.fields
}
