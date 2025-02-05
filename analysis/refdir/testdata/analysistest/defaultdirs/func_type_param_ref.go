package defaultdirs

func TestTypeParamRefs[T any](param T) T {
	var _ T
	return param
}
