package helpers

func GetPrt[T any](v any) *T {
	prt := new(T)

	*prt = v
	return prt
}

func GetPtrOrNil[T any](v any) *T {
	if v == nil {
		return nil
	}

	return new(T)
}
