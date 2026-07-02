# vmware_licenses_exporter

[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)

A **VMware vSphere** license exporter for the Prometheus/Grafana stack, built on
[`github.com/fjacquet/licenses-exporter-core`](https://github.com/fjacquet/licenses-exporter-core).
It periodically polls the vCenter `LicenseManager` for one or more vCenters and normalizes
capacity/usage/expiration data into the shared `license_` Prometheus schema, exposed via
**both** a Prometheus `/metrics` endpoint and an OTLP metric push, fed from a single shared
snapshot.

Part of the `licenses_exporter` family; shares the `license_` schema via
`licenses-exporter-core` — see [ADR-0001](docs/adr/0001-consume-core-retain-govmomi.md).

## Metrics

One `license_` prefix shared across the family; vendors are distinguished by labels, not by
metric name:

| Metric | Labels | Notes |
|---|---|---|
| `license_seats_total` | `vendor,product,unit,instance` | Omitted for unlimited license keys (`Total <= 0`) — never a `0`/`9999` sentinel. |
| `license_seats_used` | `vendor,product,unit,instance` | Raw fact, always emitted when known. |
| `license_expiration_timestamp_seconds` | `vendor,product,instance` | Omitted entirely for perpetual licenses. |
| `license_up` | `vendor,instance` | `1`/`0` per source's last collection cycle. |
| `license_collector_last_success_timestamp_seconds` | `vendor,instance` | Unix timestamp of the last successful collection. |
| `license_scrape_duration_seconds` | `vendor,instance` | Time spent collecting that source. |
| `license_build_info` | `version,goversion` | Constant `1`; exporter build metadata. |

No exporter-computed `days_to_expiration` or compliance verdict — derive those in PromQL /
alert rules from the raw facts above. An unparseable value yields an absent sample, never a
fake `0`. At cold start only `license_build_info` is emitted; per-target series appear once
each vCenter's first collection cycle resolves. See [docs/metrics.md](docs/metrics.md) for
the full reference.

## Quick start

```bash
make cli
./bin/vmware_licenses_exporter --config config.yaml
# metrics: http://localhost:9106/metrics   health: http://localhost:9106/health
```

Useful flags: `--once --debug` runs a single collection cycle and dumps every collected
sample (sorted, exposition style) instead of serving; `--trace` logs repo-owned API response
bodies for live payload validation (never SDK debug modes, which would leak the session
cookie).

## Configuration

The VMware collector is toggled in `config.yaml` (`vmware.enabled`), not via environment
variables. Secrets are referenced as `${ENV}` placeholders inside `config.yaml` (or via
`passwordFile` for file-based secrets); a `.env` file is a convenience for local `${ENV}`
expansion, never the source of truth. See `config.yaml` for a full example (one or more
vCenters).

### vCenter read-only role

The collector authenticates as a configured vCenter account, calls `LicenseManager.List`,
then logs out — stateless, once per collection cycle. Use a dedicated, read-only account
granted only the `Global.Licenses` privilege (to list license keys) plus the minimal
`Sessions` privileges needed to log in and out. No inventory, VM, or configuration
privileges are required.

## Demo stack

```bash
docker compose up
```

Brings up the exporter (`:9106`), Prometheus, and Grafana, auto-provisioned. See
[docs/deployment/docker.md](docs/deployment/docker.md) and
[docs/dashboards.md](docs/dashboards.md).

## Development

```bash
make tools   # install golangci-lint, cyclonedx-gomod, govulncheck (pinned)
make ci      # gofmt check + vet + lint + race tests + govulncheck + build (the CI gate)
```

## License

Apache License 2.0 — see [LICENSE](LICENSE).
