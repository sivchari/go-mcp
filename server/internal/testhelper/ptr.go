package testhelper

func To[T any](x T) *T {
	return &x
}
