package defaultdirs

var StringVarAtStart string

func TestRefUpToStringVar() {
	_ = StringVarAtStart // want "var reference StringVarAtStart is after definition"
}
