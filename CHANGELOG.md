# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.1.3] - 2026-09-13

### Security

- `google.golang.org/grpc` (transitive, via `licenses-exporter-core`) 1.83.0 -> 1.83.2,
  fixing **GHSA-vp52-pcj8-j9qc** and **GHSA-2v4p-qf9q-27wj** (HIGH).

### Changed

- `github.com/fjacquet/licenses-exporter-core` 1.1.1 -> 1.1.2, also carrying an `otel`
  1.45.0 -> 1.46.0 refresh and a `logrus` 1.10.0 -> 1.10.2 bump.
- `github.com/vmware/govmomi` 0.55.1 -> 0.56.0.
- Docker base image `golang` 1.26.6 -> 1.27.1.
- Dependabot auto-merge enabled for this repo, then hardened: the bot-actor guard now
  reads `github.event.pull_request.user.login` instead of the spoofable
  `github.actor`.

## [2.0.0] - 2026-08-01

### Breaking
- The published container image runs as uid **10001** (named user `licenses`),
  not `65532`. `Dockerfile.goreleaser` moves from `gcr.io/distroless/static:nonroot`
  to `alpine:latest`, matching the local `./Dockerfile` and the rest of the
  exporter family. Anyone pinning the container uid — a `securityContext.runAsUser`,
  ownership on a mounted secret or log volume — must update it. See ADR-0002.

### Security
- Bump `licenses-exporter-core` to v1.1.1, carrying its grpc v1.82.0 -> v1.83.0
  fix for GO-2026-6061, a reachable vulnerability. v1.1.1 is documentation-only
  vs v1.1.0 — no code changes on top of the bump below.

### Added
- **`/livez` and `/readyz`**, both always 200 and reading no state, via
  `licenses-exporter-core` v1.1.1. Point Kubernetes probes and container
  healthchecks at these, never at `/metrics`.
- **`HEALTHCHECK`** against `http://127.0.0.1:9106/livez` in both `Dockerfile` and
  `Dockerfile.goreleaser`, and a matching `healthcheck:` in `docker-compose.yml`
  and `docker-compose.ghcr.yml`. The ghcr healthcheck needs an image from this
  release or later — earlier published images are distroless and carry no `wget`.

### Changed
- **`/health` now always returns 200**, with `starting`/`ok` as the body rather
  than the status code. It previously returned 503 until the first collection
  cycle completed, which restarted healthy pods under a `livenessProbe` and
  reported containers unhealthy for the whole start-up window. Anything asserting
  a 503 from `/health` must be updated.
- **`/metrics` now accepts `name[]` query parameters** for metric filtering, a
  behavior change in `prometheus/client_golang` 1.24.1 (pulled in via the
  `licenses-exporter-core` bump above) that applies automatically to this
  exporter's plain `promhttp.HandlerFor(reg, promhttp.HandlerOpts{})`. Scraping
  with no `name[]` parameter is unchanged. Separately, that same release always
  validates metric names against the UTF-8 scheme; neither this exporter nor
  `licenses-exporter-core` sets `model.NameValidationScheme = LegacyValidation`,
  so there is no observable effect here.

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
