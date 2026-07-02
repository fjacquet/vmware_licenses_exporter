# 1. Consume licenses-exporter-core and retain govmomi

Date: 2026-07-02

## Status
Accepted

## Context
This exporter is the VMware sibling in the licenses_exporter family. The
vendor-neutral engine (schema, snapshot store, collection loop, dual export,
hot-reload server) lives in the shared library
`github.com/fjacquet/licenses-exporter-core`. For vSphere access, govmomi is the
official, mature Go SDK: it models the `LicenseManager` view and the license
capacity/usage/expiration fields we export, with session-based login.

## Decision
Depend on `licenses-exporter-core`; build every sample through its constructors.
`main.go` delegates the whole lifecycle to `core.Main`. Retain govmomi as the
vSphere client — it is both available and useful per the family client rule, so no
hand-rolled resty client and no "SDK-not-useful" ADR is warranted. Authenticate
stateless per collection cycle (login → LicenseManager.List → logout), with no
persisted session.

## Consequences
- Schema identity is guaranteed by construction — no local `license_` metric code.
- Engine bugfixes/features arrive via a core version bump, not a local edit.
- govmomi's dependency weight is accepted for its correctness and coverage.
- Startup is fatal on an unbuildable-but-valid config (core behaviour); see the core CHANGELOG.
