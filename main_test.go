package main

import (
	"os"
	"path/filepath"
	"testing"

	core "github.com/fjacquet/licenses-exporter-core"
)

// TestLoadConfigParsesBaseAndVMware proves the consumer Config wires core.Base
// (collection/otlp) AND the vendor vmware block from one YAML file.
func TestLoadConfigParsesBaseAndVMware(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	yaml := `
collection:
  interval: 3h
otlp:
  endpoint: "otel:4317"
  insecure: true
vmware:
  enabled: true
  vcenters:
    - instance: dc-a
      host: https://vcenter-a.example.com
      username: svc-ro
      password: shhh
      insecureSkipVerify: true
`
	if err := os.WriteFile(path, []byte(yaml), 0o600); err != nil {
		t.Fatal(err)
	}
	var cfg Config
	if err := core.LoadYAML(path, &cfg); err != nil {
		t.Fatalf("LoadYAML: %v", err)
	}
	if cfg.Collection.Interval.Hours() != 3 {
		t.Errorf("interval = %v, want 3h", cfg.Collection.Interval)
	}
	if cfg.OTLP.Endpoint != "otel:4317" {
		t.Errorf("otlp endpoint = %q, want otel:4317", cfg.OTLP.Endpoint)
	}
	if !cfg.VMware.Enabled || len(cfg.VMware.VCenters) != 1 || cfg.VMware.VCenters[0].Instance != "dc-a" {
		t.Errorf("vmware block not parsed: %+v", cfg.VMware)
	}
}

// TestLoadReturnsSourcesForEnabledVCenter proves the App.Load closure builds a
// core.Source per enabled vCenter (the wiring core.Main will drive).
func TestLoadReturnsSourcesForEnabledVCenter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	yaml := `
collection:
  interval: 2h
vmware:
  enabled: true
  vcenters:
    - instance: dc-a
      host: https://vcenter-a.example.com
      username: svc-ro
      password: shhh
`
	if err := os.WriteFile(path, []byte(yaml), 0o600); err != nil {
		t.Fatal(err)
	}
	base, sources, err := loadConfig(path)
	if err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if base.Collection.Interval.Hours() != 2 {
		t.Errorf("interval = %v, want 2h", base.Collection.Interval)
	}
	if len(sources) != 1 {
		t.Fatalf("got %d sources, want 1", len(sources))
	}
	if sources[0].Vendor() != "vmware" || sources[0].Instance() != "dc-a" {
		t.Errorf("source identity = %s/%s, want vmware/dc-a", sources[0].Vendor(), sources[0].Instance())
	}
}
