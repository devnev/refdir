package defaultdirs

type TestTypeRefUpType struct{}

func TestTypeRefUp() {
	_ = TestTypeRefUpType{} // want "type reference TestTypeRefUpType is after definition"
}
