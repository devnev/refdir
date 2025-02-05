package refdir

import (
	"os"
	"path/filepath"
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
