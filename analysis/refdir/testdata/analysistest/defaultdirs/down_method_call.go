package defaultdirs

type TestDownCallToMethod struct{}

func (s TestDownCallToMethod) TestDownRefToMethod() {
	s.DummyMethod()
}

func (s TestDownCallToMethod) DummyMethod() {}
