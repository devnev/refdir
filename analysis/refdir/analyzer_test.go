package refdir

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer_DefaultDirs(t *testing.T) {
	colorize = false
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failted to get workdir: %v", err)
	}
	analysistest.Run(t, filepath.Join(wd, "testdata/analysistest"), Analyzer, "./defaultdirs/...")
}

func TestDefaultRefOrderIsValid(t *testing.T) {
	for _, kind := range RefKinds {
		if _, ok := RefOrder[kind]; !ok {
			t.Errorf("ref kind %v missing from RefOrder", kind)
		}
	}
	for kind, dir := range RefOrder {
		if !slices.Contains(RefKinds, kind) {
			t.Errorf("invalid kind %v in RefOrder", kind)
		}
		if !slices.Contains(Directions, dir) {
			t.Errorf("invalid direction %v for kind %v in RefOrder", dir, kind)
		}
	}
}
