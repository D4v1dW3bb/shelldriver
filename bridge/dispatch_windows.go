package bridge

func Dispatch(fn func() error) error {
	return fn()
}
