package main

import (
	"github.com/devnev/refdir/analysis/refdir"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(refdir.Analyzer) }
