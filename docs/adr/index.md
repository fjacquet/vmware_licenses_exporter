# Architecture decision records

This directory records the significant architectural decisions for
`vmware_licenses_exporter` — the *why* behind the design, in the form of dated
[MADR](https://adr.github.io/madr/)-style records. Decisions are immutable once
accepted: rather than editing a past record, add a new one that supersedes it.

| ADR | Decision | Status |
|---|---|---|
| [0001](0001-consume-core-retain-govmomi.md) | Consume `licenses-exporter-core` and retain `govmomi` as the vSphere client | accepted |

To add a decision, copy [`0001`](0001-consume-core-retain-govmomi.md)'s structure to
the next number and link it here.
