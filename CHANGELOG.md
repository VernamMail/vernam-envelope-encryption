# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- **Breaking (domain separation):** the canonical HKDF info string is now the vendor-neutral `"EnvelopeMetadataEncryption-v1"`. The interim `"VernamMail-EnvelopeEncryption-v1"` from the 0.1.0 draft was never shipped by any implementation and is not recognized; the deployed legacy string remains `"ENIGMA-HybridKEM-v1"` (SPEC §4.5.1)
- Go module path corrected to `github.com/VernamMail/vernam-envelope-encryption/go` so the module resolves from this repository
- Specification corrections: §7.8 IKM description (32-byte ML-KEM-1024 shared secret; 1568 bytes is the ciphertext size), HKDF written as the full RFC 5869 extract-then-expand construction, §4.5.1 legacy-mode wording moved to future tense until milestone 3, storage-overhead arithmetic aligned to the 13 fields defined in §3
- Documentation refresh: prior-art facts corrected (TutaCrypt shipped March 2024; Proton subject line is not end-to-end encrypted; ETSI reference titles and dates fixed), funding status updated to target the first Restack open call

### Added

- Minimal standalone TypeScript reference implementation of envelope-field encryption (`ts/envelope-field.ts`), depending only on the Web Crypto API
- Cross-language interop test (`TestKnownVectorMatchesGo`) asserting the TypeScript output is byte-identical to the Go reference for test vector envelope-field-001
- `ts/README.md` documenting the cross-language reference and its scope

### Planned

The milestones toward v0.5 and v0.8 (extraction, fuzzing, expanded vectors, hybrid KEM, interop suite), planned for grant funding, are tracked in [ROADMAP.md](./ROADMAP.md).

## [0.1.0] - 2026-05-02

Initial public draft.

### Added

- Specification draft v0.1 ([SPEC.md](./SPEC.md), 10 sections, ~465 lines)
- Threat model with explicit adversary classes A1-A7 ([docs/threat-model.md](./docs/threat-model.md))
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
