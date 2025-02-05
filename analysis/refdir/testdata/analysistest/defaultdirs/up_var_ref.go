package defaultdirs

var StringVarAtStart string

func TestRefUpToStringVar() {
	_ = StringVarAtStart // want "value reference StringVarAtStart is after definition"
}
