package defaultdirs

type EmptyStruct struct{}

func (e EmptyStruct) TestReceiverReference() {
	_ = e
}
