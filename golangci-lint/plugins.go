package golangcilint

import (
	"github.com/devnev/refdir/analysis/refdir"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("refdir", New)
}

func New(settings any) (register.LinterPlugin, error) {
	return &Plugin{}, nil
}

type Plugin struct{}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{refdir.Analyzer}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
