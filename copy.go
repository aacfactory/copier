package copier

func Copy(dst any, src any) (err error) {
	err = DefaultConfig.Copy(dst, src)
	return
}

func From[D any](src any) (dst D, err error) {
	//if typ := reflect2.TypeOf(dst); typ.Kind() == reflect.Ptr {
	//	dst = typ.New()
	//	err = Copy(dst, src)
	//} else {
	//	err = Copy(&dst, src)
	//}
	return
}
