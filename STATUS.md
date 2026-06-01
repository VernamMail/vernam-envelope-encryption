# Status

**Last updated:** 2026-05-02

## Current Phase

**Specification drafted; partial reference implementation; pending NLnet Restack fund support for full extraction, hardening, audit, and v1.0 release.**

This repository represents the open-commons extraction of a protocol that is already in production use within Vernam Mail's pre-launch encrypted-email system. The grant funds the **extraction work, not the original implementation**.

## What Is Done

- **Specification draft (v0.1):** [SPEC.md](./SPEC.md) — 10 sections, ~445 lines, including threat model, wire formats, and compliance criteria
- **Threat model:** [docs/threat-model.md](./docs/threat-model.md) — explicit adversary classes A1–A7, trust boundaries, known limitations
- **Wire format documentation:** [docs/wire-format.md](./docs/wire-format.md) — byte-level diagrams
- **Prior art comparison:** [docs/prior-art.md](./docs/prior-art.md) — situated against Proton, Tuta, Sequoia PGP, Rosenpass, oqsprovider
- **NGI alignment:** [docs/ngi-alignment.md](./docs/ngi-alignment.md) — how this project advances the Next Generation Internet vision
- **Reference implementation, partial:** Go library at [go/](./go/) provides:
  - 32-byte session key generation (`NewSessionKey`)
  - AES-256-GCM envelope-field encryption / decryption (`EncryptField`, `DecryptField`)
  - Wire-format constants validated against the specification
  - 10 passing tests including a verified known-vector test against [test-vectors/basic.json](./test-vectors/basic.json)
- **Verified test vector:** envelope-field-001 in [test-vectors/basic.json](./test-vectors/basic.json), reproducible via `go test -v ./...`
- **License:** [Apache 2.0](./LICENSE)
- **Vulnerability disclosure policy:** [SECURITY.md](./SECURITY.md)
- **Contribution process:** [CONTRIBUTING.md](./CONTRIBUTING.md)

## What Is Not Done (Grant-Funded Scope)

These are the deliverables proposed in [ROADMAP.md](./ROADMAP.md):

- **Specification finalization** through a public 4-week review period and revision
- **Production-grade reference library extraction**: hardening the envelope-field encryption beyond the current skeleton (fuzzing, comprehensive test vectors, godoc, integration patterns, decoupled types)
- **Hybrid KEM library extraction** from product code: ML-KEM-1024 + X25519 with v2 ciphertext-binding AAD
- **Cross-language interop validation:** documented test vectors and test harness proving the Go reference library, the existing TypeScript implementation, and the existing Kotlin implementation produce and consume identical wire formats
- **Migration path** from the legacy HKDF info string (see SPEC.md §4.5.1) to the canonical name
- **Tutorial, examples, and integration guide**
- **Third-party cryptographic review** by an independent academic or professional firm
- **Tagged v1.0 release** with the cryptographic-review report attached

## Relationship to Existing Implementations

The protocol described in this repository is currently implemented inside Vernam Mail's private, product-coupled codebase across Go (server), TypeScript (web client), and Kotlin (Android). These implementations have been interoperating with each other in production-grade conditions for ~5 months.

The grant **does not fund** the existing private implementations. The grant funds the work to:

1. Decouple the protocol from product-specific data structures, error models, and database schemas
2. Produce a clean, vendor-neutral public Go API
3. Generate a comprehensive, publicly-auditable test vector suite
4. Document the protocol with sufficient rigor that implementers in other languages and other systems can adopt it
5. Commission and respond to a third-party cryptographic review
6. Publish a stable v1.0 release independent of any product

These artifacts do not exist today and would not exist without grant funding. Vernam Mail's pre-launch budget covers product work, not commons work.

## Funding Status

NLnet Restack fund application target: deadline 2026-08-01.

This repository will be updated upon decision. If the application is not approved, the specification and partial reference implementation will continue to evolve as time permits, but the milestones in [ROADMAP.md](./ROADMAP.md) are predicated on grant funding.

## Contact

- **Specification questions, eligibility, or interop inquiries:** open a GitHub issue (or email contact below until issues are enabled)
- **Security reports:** see [SECURITY.md](./SECURITY.md)
- **General correspondence:** `security@vernammail.com`
