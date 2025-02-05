package defaultdirs

const StringConstAtStart = "top"

func TestRefUpToStringConst() {
	_ = StringConstAtStart // want "const reference StringConstAtStart is after definition"
}
