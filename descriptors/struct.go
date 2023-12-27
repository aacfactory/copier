package descriptors

import (
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
	"unsafe"
)

type StructFieldDescriptor struct {
	Name  string
	Field []*reflect2.UnsafeStructField
	anon  *reflect2.UnsafeStructType
}

func (desc *StructFieldDescriptor) ValueOf(ptr unsafe.Pointer) (v unsafe.Pointer, typ reflect2.Type, err error) {
	fieldPtr := ptr
	for i, f := range desc.Field {
		fieldPtr = f.UnsafeGet(fieldPtr)

		if i == len(desc.Field)-1 {
			typ = f.Type()
			break
		}

		if f.Type().Kind() == reflect.Ptr {
			if *((*unsafe.Pointer)(fieldPtr)) == nil {
				err = fmt.Errorf("copier: embedded field %q is nil", f.Name())
				return
			}
			fieldPtr = *((*unsafe.Pointer)(fieldPtr))
		}
		typ = f.Type()
	}
	v = fieldPtr
	return
}

type StructFields []*StructFieldDescriptor

func (sf StructFields) Get(name string) *StructFieldDescriptor {
	for _, f := range sf {
		if f.Name == name {
			return f
		}
	}
	return nil
}

func (sf StructFields) GetByTag(tagKey string, name string) *StructFieldDescriptor {
	for _, f := range sf {
		for _, field := range f.Field {
			if tag, has := field.Tag().Lookup(tagKey); has && tag == name {
				return f
			}
		}
	}
	return nil
}

func DescribeStruct(typ reflect2.Type) *StructDescriptor {
	structType := typ.(*reflect2.UnsafeStructType)
	fields := StructFields{}

	var curr []StructFieldDescriptor
	next := []StructFieldDescriptor{{anon: structType}}

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

					next = append(next, StructFieldDescriptor{Field: chain, anon: t.(*reflect2.UnsafeStructType)})
					continue
				}

				// Ignore unexported fields.
				if isUnexported {
					continue
				}

				fieldName := field.Name()

				fields = append(fields, &StructFieldDescriptor{
					Name:  fieldName,
					Field: chain,
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

func (desc *StructDescriptor) Field(name string) (f *StructFieldDescriptor) {
	f = desc.fields.Get(name)
	return
}

func (desc *StructDescriptor) FieldByTag(tagKey string, name string) (f *StructFieldDescriptor) {
	f = desc.fields.GetByTag(tagKey, name)
	return
}
