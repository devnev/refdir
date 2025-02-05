package defaultdirs

type TestRefUpConstraintType interface {
	DummyMethod()
}

func TestRefUpToConstraintType[T TestRefUpConstraintType]() {} // want "type reference TestRefUpConstraintType is after definition"
