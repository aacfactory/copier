package descriptors

import (
	"github.com/modern-go/reflect2"
	"time"
)

type NativeDescriptor struct {
	typ reflect2.Type
}

func (desc *NativeDescriptor) Type() reflect2.Type {
	return desc.typ
}

func DescribeString() Descriptor {
	return &NativeDescriptor{
		typ: reflect2.TypeOf(""),
	}
}

func DescribeBool() Descriptor {
	return &NativeDescriptor{
		typ: reflect2.TypeOf(false),
	}
}

func DescribeInt() Descriptor {
	return &NativeDescriptor{
		typ: reflect2.TypeOf(int64(0)),
	}
}

func DescribeFloat() Descriptor {
	return &NativeDescriptor{
		typ: reflect2.TypeOf(float64(0)),
	}
}

func DescribeUint() Descriptor {
	return &NativeDescriptor{
		typ: reflect2.TypeOf(uint64(0)),
	}
}

func DescribeTime() Descriptor {
	return &NativeDescriptor{
		typ: reflect2.TypeOf(time.Time{}),
	}
}

func DescribeBytes() Descriptor {
	return &NativeDescriptor{
		typ: reflect2.TypeOf([]byte{}),
	}
}
