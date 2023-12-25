package writers

import (
	"github.com/modern-go/reflect2"
	"unsafe"
)

type Reader interface {
	Type() (typ reflect2.Type)
	Read() (v unsafe.Pointer)
}
