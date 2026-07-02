package vmware

// VMwareConfig is the vSphere block of the exporter config. Enabled=false (or an
// empty VCenters list) yields zero sources — the exporter then serves only
// license_build_info.
type VMwareConfig struct {
	Enabled  bool            `yaml:"enabled"`
	VCenters []VCenterConfig `yaml:"vcenters"`
}

// VCenterConfig is one vCenter target. Password is an inline ${ENV} ref;
// PasswordFile is a path read at load (ResolveSecret governs precedence).
type VCenterConfig struct {
	Instance           string `yaml:"instance"`
	Host               string `yaml:"host"`
	Username           string `yaml:"username"`
	Password           string `yaml:"password"`
	PasswordFile       string `yaml:"passwordFile"`
	InsecureSkipVerify bool   `yaml:"insecureSkipVerify"`
}
