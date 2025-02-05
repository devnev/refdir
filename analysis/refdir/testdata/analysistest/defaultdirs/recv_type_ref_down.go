package defaultdirs

func (s TestDownRecvTypeRef) TestDownRecvTypeRef() { // want "recvtype reference TestDownRecvTypeRef is before definition"
	_ = TestDownRecvTypeRef{}
}

type TestDownRecvTypeRef struct{}
