package golangcilint

import (
	"fmt"
	"slices"

	"github.com/devnev/refdir/analysis/refdir"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("refdir", New)
}

func New(rawSettings any) (register.LinterPlugin, error) {
	settings, err := register.DecodeSettings[map[string]string](rawSettings)
	if err != nil {
		return nil, err
	}

	parsedOrder := map[refdir.RefKind]refdir.Direction{}
	for key, value := range settings {
		if !slices.Contains(refdir.RefKinds, refdir.RefKind(key)) {
			return nil, fmt.Errorf("invalid refdir settings key %q", key)
		}
		if !slices.Contains(refdir.Directions, refdir.Direction(value)) {
			return nil, fmt.Errorf("invalid refdir direction %q for settings key %q", value, key)
		}
		parsedOrder[refdir.RefKind(key)] = refdir.Direction(value)
	}

	return &Plugin{settings: parsedOrder}, nil
}

type Plugin struct {
	settings map[refdir.RefKind]refdir.Direction
}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	for kind, dir := range p.settings {
		refdir.RefOrder[kind] = dir
	}
	if err := refdir.Analyzer.Flags.Set("color", "false"); err != nil {
		return nil, fmt.Errorf("failed to disable color setting: %w", err)
	}
	return []*analysis.Analyzer{refdir.Analyzer}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
