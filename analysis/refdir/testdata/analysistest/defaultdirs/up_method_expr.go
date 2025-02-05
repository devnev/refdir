package defaultdirs

type TestUpRefToMethod struct{}

func (s TestUpRefToMethod) DummyMethod() {}

func (s TestUpRefToMethod) TestDownRefToMethod() {
	_ = s.DummyMethod // want "func reference DummyMethod is after definition"
}
