package main

import (
	core "github.com/fjacquet/licenses-exporter-core"
	"github.com/fjacquet/vmware_licenses_exporter/internal/vmware"
)

// Config is the exporter's full config: the shared core.Base (collection + otlp)
// inline, plus the vendor-specific vmware block.
type Config struct {
	core.Base `yaml:",inline"`
	VMware    vmware.VMwareConfig `yaml:"vmware"`
}

// loadConfig parses the file and builds the sources — the single closure body
// core.Main calls at startup and on every reload.
func loadConfig(path string) (core.Base, []core.Source, error) {
	var cfg Config
	if err := core.LoadYAML(path, &cfg); err != nil {
		return core.Base{}, nil, err
	}
	if err := cfg.Validate(); err != nil {
		return core.Base{}, nil, err
	}
	sources, err := vmware.NewSources(cfg.VMware)
	if err != nil {
		return core.Base{}, nil, err
	}
	return cfg.Base, sources, nil
}
