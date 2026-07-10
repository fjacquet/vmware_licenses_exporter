# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.2] - 2026-07-10

### Security
- Bump Go to 1.26.5 to patch GO-2026-5856 (crypto/tls), reported by govulncheck.

### Added
- First multi-arch (linux/amd64, linux/arm64) container image, published to
  `ghcr.io/fjacquet/vmware_licenses_exporter` by the release pipeline via GoReleaser
  `dockers_v2`.

### Added
- Initial release: a VMware vSphere license exporter (vCenter `LicenseManager` via govmomi)
  built on `github.com/fjacquet/licenses-exporter-core`. Emits the shared `license_` schema
  (`vendor="vmware"`), so it shares one Prometheus / Grafana view with the other family
  exporters. Default metrics port `9106`. See ADR-0001.
