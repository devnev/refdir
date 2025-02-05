package defaultdirs

func TestRefDownToConstraintType[T TestRefDownConstraintType]() {}

type TestRefDownConstraintType interface {
	DummyMethod()
}
