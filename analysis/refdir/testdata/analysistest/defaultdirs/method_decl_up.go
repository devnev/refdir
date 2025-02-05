package defaultdirs

func (s TestUpDeclOfMethod) TestUpDeclOfMethod() {} // want "recvtype reference TestUpDeclOfMethod is before definition"

type TestUpDeclOfMethod struct{}
