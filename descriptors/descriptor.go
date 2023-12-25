package descriptors

import "github.com/modern-go/reflect2"

type Descriptor interface {
	Type() reflect2.Type
}
