# Metrics reference

`vmware_licenses_exporter` exposes one generic `license_` metric family, shared across the
`licenses_exporter` family via `github.com/fjacquet/licenses-exporter-core`. Vendors are
distinguished by **labels**, not by metric name — this exporter emits `vendor="vmware"`.
Every value is a raw fact straight from the vSphere `LicenseManager`: there is no
exporter-computed compliance verdict or "days remaining" gauge. Derive those in PromQL or
alert rules from the raw facts below.

This table is the diff target for `--once --debug`, which dumps every collected sample
(sorted, exposition style) for live payload validation against a real vCenter.

## Example series

```text
license_seats_total{vendor="vmware",product="vSphere 8 Enterprise Plus",unit="CPUs",instance="dc-a"} 32
license_seats_used{vendor="vmware",product="vSphere 8 Enterprise Plus",unit="CPUs",instance="dc-a"} 24
license_expiration_timestamp_seconds{vendor="vmware",product="vSphere 8 Enterprise Plus",instance="dc-a"} 1.8039456e+09
```

## License facts

| Metric | Type | Labels | Meaning |
|---|---|---|---|
| `license_seats_total` | Gauge | `vendor, product, unit, instance` | Total license capacity (`LicenseManagerLicenseInfo.Total`). **Omitted** when the key is unlimited (`Total <= 0`) — never a `0` or `9999` sentinel. |
| `license_seats_used` | Gauge | `vendor, product, unit, instance` | Currently consumed capacity (`LicenseManagerLicenseInfo.Used`). Always emitted when known — an unlimited key still reports its `Used` count. |
| `license_expiration_timestamp_seconds` | Gauge | `vendor, product, instance` | License expiration as a Unix timestamp, parsed from the key's `expirationDate` property. **Omitted entirely** when the license is perpetual (no `expirationDate` property present). |

## Health / state

| Metric | Type | Labels | Meaning |
|---|---|---|---|
| `license_up` | Gauge | `vendor, instance` | `1` if the vCenter's last collection cycle (login → `LicenseManager.List` → logout) succeeded, `0` if it failed. Absent entirely until that vCenter's first cycle resolves. |
| `license_collector_last_success_timestamp_seconds` | Gauge | `vendor, instance` | Unix timestamp of the last successful collection for this vCenter. `time() - this` is the data-age/freshness signal. |
| `license_scrape_duration_seconds` | Gauge | `vendor, instance` | Wall-clock time spent collecting this vCenter during the last cycle. |
| `license_build_info` | Gauge | `version, goversion` | Constant `1`; carries the exporter's build metadata. The only series present before the first collection cycle completes. |

## Label semantics

| Label | Meaning / source |
|---|---|
| `vendor` | `"vmware"`. |
| `product` | The license key's `Name` (e.g. `vSphere 8 Enterprise Plus`). Raw vendor identifier — no friendly-name mapping. |
| `unit` | The key's `CostUnit` (e.g. `CPUs`); falls back to `unit` if vCenter reports an empty cost unit. |
| `instance` | The configured vCenter's `instance` name from `config.yaml` (e.g. `dc-a`). One process can poll many vCenters. |

## Design rules (raw facts, absent-not-zero)

- **Unlimited keys omit `seats_total`.** A license key with `Total <= 0` (vSphere's convention
  for "unlimited") emits only `license_seats_used` — never a `0`, and never a large sentinel
  standing in for "unlimited".
- **Perpetual licenses omit expiration.** A key with no `expirationDate` property never emits
  `license_expiration_timestamp_seconds` — there is no `9999`-year row to filter out in
  dashboards or alerts.
- **No `days_to_expiration` gauge, no perpetual sentinel.** `license_expiration_timestamp_seconds`
  carries the absolute Unix timestamp; compute days remaining in PromQL:
  `(license_expiration_timestamp_seconds - time()) / 86400`.
- **No exporter-computed `compliance_status`.** Over-allocation is
  `license_seats_used > license_seats_total`; policy belongs in PromQL/alert rules, not the
  exporter.
- **Absent, never zero.** An unparseable or missing capacity/used value yields an *absent*
  sample, never a fake `0` — a false `0` on a capacity metric would silently corrupt
  dashboards and over-allocation alerts.
- **Cold start.** Immediately after startup, before any vCenter's first collection cycle
  resolves, `/metrics` exposes **only** `license_build_info` — no `license_up` or per-target
  series exist yet, so a scrape during that window can never see a transient `0` or a
  flapping target.
- **Label-key consistency.** Every series of a given metric name carries the same label-key
  set, built from the shared constructors in `licenses-exporter-core` (see
  [ADR-0001](adr/0001-consume-core-retain-govmomi.md)).

## Live validation

```bash
./bin/vmware_licenses_exporter --config config.yaml --once --debug
```

Runs a single collection cycle and prints every collected sample in sorted, Prometheus
exposition-style output — diff it against the tables above to catch a silently-absent
metric that `license_up` alone would not reveal.
