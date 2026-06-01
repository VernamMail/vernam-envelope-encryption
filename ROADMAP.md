# Roadmap

This roadmap aligns with the deliverables proposed to NLnet's Restack fund. Milestones reflect the work required to take the existing private, product-coupled protocol implementation and produce a clean, audited, public commons artifact.

The current public repository contains an initial draft specification and a skeletal Go reference implementation of envelope-field encryption sufficient to demonstrate technical capability and reproducible test vectors. Grant funding supports the substantial extraction, formalization, and validation work that remains.

## Milestone 1 — Specification Finalization (Months 1–2)

**State:**
- [x] Initial draft specification (v0.1) at [SPEC.md](./SPEC.md)
- [x] Threat model with explicit adversary classes
- [x] Wire-format byte-level documentation
- [x] Prior-art comparison
- [ ] Public review period (4 weeks)
- [ ] External cryptographic review of the draft specification (informal, pre-implementation)
- [ ] Incorporation of community feedback
- [ ] Resolution of legacy-info-string migration path (see SPEC.md §4.5.1)
- [ ] Tagged release: SPEC v1.0 (Markdown + rendered HTML)

**Deliverable:** finalized SPEC.md, threat-model.md, wire-format.md tagged at v1.0; public review log archived.

## Milestone 2 — Reference Library Extraction & Hardening (Months 3–5)

The current public skeleton implements envelope-field encryption at v0.1. This milestone extracts the production-grade implementation from product-coupled code and hardens the public library to release quality.

**Tasks:**
- [ ] Extract envelope-field encryption from Vernam Mail's product code into the public library, decoupled from product-specific types, error models, and database schemas
- [ ] Add per-recipient session-key wrapping primitives (the symmetric and PGP-fallback paths; not yet the hybrid KEM)
- [ ] Comprehensive unit tests across boundaries: empty plaintext, large plaintext, malformed inputs, tampering, version mismatch
- [ ] Fuzzing harness for wire-format parsing (`go test -fuzz`)
- [ ] Expand test vector suite (at least 10 envelope-field vectors covering edge cases)
- [ ] godoc-quality documentation on every exported function
- [ ] Integration test suite usable by implementers in other languages
- [ ] Module published to Go module proxy

**Deliverable:** Go module v0.5, published, with documentation, tests, fuzzing, and expanded test vectors.

## Milestone 3 — Hybrid KEM Extraction & Cross-Language Interop (Months 6–8)

**Tasks:**
- [ ] Extract the hybrid ML-KEM-1024 + X25519 wrap/unwrap implementation from product code into the public Go library, decoupled
- [ ] Implement v2 ciphertext-binding AAD per SPEC §4.3
- [ ] Implement v1 backwards-compatibility decryption path
- [ ] Implement legacy-info-string migration mode per SPEC §4.5.1 (decryption fallback for production-deployed legacy ciphertexts)
- [ ] Cross-language interop test suite: documented test vectors that the Go reference, the existing TypeScript implementation, and the existing Kotlin implementation must each produce and consume identically
- [ ] Bug fixes from interop testing (any divergences found between the three implementations are resolved by updating implementations to match the spec)
- [ ] Public report: "Interop validation across three independent implementations of the envelope-encryption protocol"

**Deliverable:** Go module v0.8, hybrid KEM working, cross-language interop validated and documented.

## Milestone 4 — Documentation, Examples, and Integration Guide (Months 9–10)

**Tasks:**
- [ ] Tutorial: "Adding envelope-metadata encryption to a mail system" (~3000 words)
- [ ] Migration guide: "Moving from OpenPGP-only to envelope encryption"
- [ ] 3–5 working example programs in `go/examples/`
- [ ] Hosted documentation site (GitHub Pages or similar)
- [ ] Comparison and positioning document: how this complements OpenPGP, S/MIME, draft-ietf-openpgp-pqc
- [ ] Public review of beta release; incorporation of feedback

**Deliverable:** documentation site, example programs, public beta announcement, feedback log.

## Milestone 5 — Third-Party Cryptographic Review and v1.0 Release (Months 11–12)

**Tasks:**
- [ ] Engage independent cryptographic reviewer (academic researcher or professional firm; candidates identified during M1)
- [ ] Review covers: hybrid KEM construction soundness, AAD-binding correctness, wire format security, side channels, library implementation quality
- [ ] Address review findings (or document as "won't fix" with rationale)
- [ ] Publish review report publicly
- [ ] Tag and release v1.0
- [ ] Public announcement (Hacker News, NLnet network, IETF mailing list pointer, blog posts)

**Deliverable:** v1.0 release with attached cryptographic review report.

## Post-v1.0 (Out of Grant Scope)

Directions beyond the funded period, mentioned for context:

- IETF Internet-Draft submission for the protocol
- Reference implementations in additional languages (Rust, TypeScript published as separate package)
- Integration shims for popular MTAs (Postfix, Haraka, etc.)
- Formal security proof in a computational model
- Continued community engagement and version maintenance

## Status Tracking

Issue tracker and progress reports will be maintained on the GitHub repository once the application result is known. Public progress reports will be published quarterly during the funded period.
