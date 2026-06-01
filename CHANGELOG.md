# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

Items planned for v0.5 (NLnet milestone 2):

- Production-grade extraction of envelope-field encryption from product code
- Comprehensive test vector suite (≥10 envelope-field vectors)
- Fuzzing harness for wire-format parsing
- Per-recipient session-key wrapping primitives (symmetric + PGP-fallback paths)
- godoc-quality documentation on all exported functions

Items planned for v0.8 (NLnet milestone 3):

- Hybrid ML-KEM-1024 + X25519 wrap/unwrap implementation
- v2 ciphertext-binding AAD per SPEC §4.3
- v1 backwards-compatibility decryption path
- Legacy-info-string migration mode per SPEC §4.5.1
- Cross-language interop test suite
- Public report on three-implementation interop validation

## [Unreleased] - cross-language reference

### Added

- Minimal standalone TypeScript reference implementation of envelope-field encryption (`ts/envelope-field.ts`), depending only on the Web Crypto API
- Cross-language interop test (`TestKnownVectorMatchesGo`) asserting the TypeScript output is byte-identical to the Go reference for test vector envelope-field-001
- `ts/README.md` documenting the cross-language reference and its scope

## [0.1.0] - 2026-05-02

Initial public draft.

### Added

- Specification draft v0.1 ([SPEC.md](./SPEC.md), 10 sections, ~465 lines)
- Threat model with explicit adversary classes A1–A7 ([docs/threat-model.md](./docs/threat-model.md))
- Wire-format documentation with byte-level diagrams ([docs/wire-format.md](./docs/wire-format.md))
- Prior-art comparison with Proton Mail, Tuta, and NLnet-funded PQC projects ([docs/prior-art.md](./docs/prior-art.md))
- NGI alignment / European Dimension document ([docs/ngi-alignment.md](./docs/ngi-alignment.md))
- Reference Go library skeleton with envelope-field encryption (AES-256-GCM)
- 10 passing tests, including a verified known-vector test against `test-vectors/basic.json`
- Test vector envelope-field-001 (verified, reproducible)
- Apache 2.0 license
- Vulnerability disclosure policy ([SECURITY.md](./SECURITY.md))
- Contribution guidelines ([CONTRIBUTING.md](./CONTRIBUTING.md))
- Code of Conduct (Contributor Covenant 2.1)
- This changelog
