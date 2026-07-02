# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release: a VMware vSphere license exporter (vCenter `LicenseManager` via govmomi)
  built on `github.com/fjacquet/licenses-exporter-core`. Emits the shared `license_` schema
  (`vendor="vmware"`), so it shares one Prometheus / Grafana view with the other family
  exporters. Default metrics port `9106`. See ADR-0001.
