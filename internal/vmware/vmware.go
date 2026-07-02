package vmware

import (
	"fmt"

	core "github.com/fjacquet/licenses-exporter-core"
)

// NewSources builds one stateless Source per configured vCenter.
func NewSources(cfg VMwareConfig) ([]core.Source, error) {
	if !cfg.Enabled {
		return nil, nil
	}
	var out []core.Source
	for _, v := range cfg.VCenters {
		pw, err := core.ResolveSecret(v.Password, v.PasswordFile)
		if err != nil {
			return nil, fmt.Errorf("vcenter %q: %w", v.Instance, err)
		}
		out = append(out, &source{
			instance: v.Instance,
			host:     v.Host,
			username: v.Username,
			password: pw,
			insecure: v.InsecureSkipVerify,
		})
	}
	return out, nil
}
