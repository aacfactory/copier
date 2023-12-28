package copier

func Copy(dst any, src any) (err error) {
	err = DefaultConfig.Copy(dst, src)
	return
}
