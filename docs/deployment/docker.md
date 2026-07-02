# Docker deployment

The image (`Dockerfile`) is a non-root, multi-stage Alpine build: it runs as the
unprivileged `licenses` user (uid `10001`), listens on `9106`, and reads `config.yaml` from
`/etc/vmware_licenses_exporter/config.yaml`.

## Standalone container

```bash
docker run -d --name vmware_licenses_exporter -p 9106:9106 \
  -e VC_HOST=https://vcenter.example.com \
  -e VC_USERNAME=svc-ro \
  -e VC_PASSWORD=... \
  -v /path/to/config.yaml:/etc/vmware_licenses_exporter/config.yaml:ro \
  ghcr.io/fjacquet/vmware_licenses_exporter:latest
```

`config.yaml` is the source of truth for which `${ENV}` references are actually consumed
(`${VC_HOST}`, `${VC_USERNAME}`, `${VC_PASSWORD}`, etc.) — every variable the mounted config
references must exist in the container's environment, or the exporter fails fast at load
with `config references unset environment variable "..."`. Secrets can alternatively be
supplied as a file via `passwordFile` in `config.yaml`, mounted as a read-only volume
instead of passed as an env var.

`/metrics` and `/health` are both served on `9106`; `/health` returns HTTP 200 with
`starting` until the first collection cycle completes for every enabled vCenter, then `ok`.

## One-command demo stack (Compose)

```bash
docker compose up
```

`docker-compose.yml` builds the exporter from the local `Dockerfile` and brings up:

- **`vmware_licenses_exporter`** (`:9106`) — built locally, config mounted from `./config.yaml`.
- **`prometheus`** (`:9090`) — scrapes the exporter per `prometheus.yml` and loads the
  alert rules in `deploy/prometheus/license.rules.yml`.
- **`grafana`** (`:3000`, `admin`/`admin` by default) — auto-provisioned with the Prometheus
  datasource and the **Enterprise Licenses — Overview** dashboard
  (`grafana/dashboards/licenses-overview.json`); see [Dashboards](../dashboards.md).

The bundled `config.yaml` ships with placeholder `${VC_*}` env references;
`docker-compose.yml` supplies default placeholder values for those so the stack starts
without any `.env` file, purely to demonstrate the wiring end-to-end. Override them (shell
env or a `.env` file next to `docker-compose.yml`) with real vCenter credentials to point
the demo at a live environment.

To run the **published** image instead of building locally:

```bash
docker compose -f docker-compose.ghcr.yml up -d
```

Pin a version with `VMWARE_LICENSES_EXPORTER_TAG` (defaults to `:latest`):

```bash
VMWARE_LICENSES_EXPORTER_TAG=0.2.1 docker compose -f docker-compose.ghcr.yml up -d
```

## Required permissions before first run

### vCenter — read-only role

The collector authenticates as the configured `username`/`password`, calls
`LicenseManager.List`, then logs out — once per collection cycle, with no persisted
session. Create a dedicated, read-only vCenter account and assign it a custom role granting
only the **`Global.Licenses`** privilege (to list/read license keys) and **`Sessions`**
privileges sufficient to log in and out (vCenter grants a minimal session automatically on
login for any authenticated account). No inventory, VM, or configuration privileges are
required. Without the `Global.Licenses` privilege, `Collect` fails with an authorization
error and that vCenter's cycle degrades to `license_up{vendor="vmware",...}=0` rather than
blocking the whole exporter.

## Flags

| Flag | Default | Meaning |
|---|---|---|
| `--config` | `config.yaml` | Path to the config file. |
| `--web.listen-address` | `:9106` | Address the HTTP server (metrics + health) binds to. |
| `--once` | `false` | Run a single collection cycle and exit instead of serving. |
| `--debug` | `false` | Debug-level logging; combined with `--once` it dumps every collected sample (sorted, exposition style) — see `docs/metrics.md`. |
| `--trace` | `false` | Logs repo-owned API responses for live payload validation. govmomi is non-injectable, so this **never** enables SDK-level debug output, which would leak the session cookie / bearer token. |

Config reload is live: `SIGHUP`, or any write/create to the config file, triggers a
validated hot reload without a restart or any interruption to `/metrics`.
