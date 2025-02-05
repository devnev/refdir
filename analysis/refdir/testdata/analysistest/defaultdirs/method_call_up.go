package defaultdirs

type TestUpCallToMethod struct{}

func (s TestUpCallToMethod) DummyMethod() {}

func (s TestUpCallToMethod) TestDownCallToMethod() {
	s.DummyMethod() // want "func reference DummyMethod is after definition"
}
