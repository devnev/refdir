package defaultdirs

type TestDownRefToMethod struct{}

func (s TestDownRefToMethod) TestDownRefToMethod() {
	_ = TestDownRefToMethod.DummyMethod
}

func (s TestDownRefToMethod) DummyMethod() {}
