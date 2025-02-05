package defaultdirs

type TestStructTypeParam[T any] struct {
	DummyField T
}

func (t TestStructTypeParam[T]) TestTypeParamRefs(param T) T {
	var _ T
	return param
}
