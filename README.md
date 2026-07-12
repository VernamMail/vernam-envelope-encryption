# vernam-envelope-encryption

> **Encrypted email systems encrypt your message body. They do not encrypt who you sent it to.** This is an open protocol, with reference implementations, that does.

## The problem

Proton Mail, Tuta, and every other major encrypted-email provider store envelope metadata (sender address, recipient addresses, CC/BCC, Reply-To, Message-ID, original date) in plaintext at rest. This is the primary input to traffic analysis: who talks to whom, and when. Intelligence disclosures and surveillance research repeatedly show that metadata, not message content, is what bulk-collection systems rely on.

## The protocol

This repository specifies and implements:

1. **Per-email symmetric encryption** of every envelope field (From, To, CC, BCC, Reply-To, Message-ID, timestamps, Subject, etc.) using AES-256-GCM with a fresh per-email session key.
2. **Per-recipient session key wrapping** using a hybrid post-quantum KEM combining ML-KEM-1024 (NIST FIPS 203) and X25519 (RFC 7748), bound together via HKDF-SHA256 with a v2 ciphertext-binding AAD construction.
3. **Wire formats** suitable for storage in standard databases and for transmission over HTTP / SMTP-adjacent protocols.
4. An explicit **threat model** that names what this defends against and what it does not.

The result: a server storing encrypted email can see the *number* of recipients and the *approximate size* of envelope fields. It cannot see who the email is from, who it is to, what its subject is, when it was originally sent, or its body.

## Status

**Specification plus partial reference implementation; pending NLnet Restack funding.** See [STATUS.md](./STATUS.md).

This repository contains:

- [SPEC.md](./SPEC.md): the protocol specification (~470 lines, 10 sections + appendix)
- [docs/threat-model.md](./docs/threat-model.md): adversary classes A1-A7, trust boundaries, known limitations
- [docs/wire-format.md](./docs/wire-format.md): byte-level diagrams
- [docs/prior-art.md](./docs/prior-art.md): comparison with Proton, Tuta, Sequoia PGP, Rosenpass, oqsprovider
- [docs/ngi-alignment.md](./docs/ngi-alignment.md): how this advances the Next Generation Internet vision
- [go/](./go/): reference implementation in Go (envelope-field encryption working; hybrid KEM scheduled for milestone 3)
- [ts/](./ts/): minimal TypeScript reference for envelope-field encryption; produces byte-identical output to the Go reference (cross-language interop demonstrated, not just claimed)
- [test-vectors/](./test-vectors/): reproducible test vectors verified by both the Go and TypeScript implementations
- [ROADMAP.md](./ROADMAP.md): five milestones aligned with NLnet deliverables
- [SECURITY.md](./SECURITY.md): vulnerability disclosure policy
- [CONTRIBUTING.md](./CONTRIBUTING.md): how to contribute during specification phase

The protocol runs in Vernam Mail's pre-launch encrypted-email deployment, where Go (server), TypeScript (web client), and Kotlin (Android) implementations have been interoperating since early 2026. Funding from NLnet (Restack fund, application planned for its first open call) will support extracting it as a standalone, audited, open-source library.

## What this is not

- Not an MTA (does not implement SMTP)
- Not a key-management system (assumes recipient public keys are known)
- Not a replacement for OpenPGP (interoperates with it)
- Not a drop-in for existing mail systems (requires schema and protocol integration)

## Reading order

- **5 minutes**: this README + [SPEC.md §1-2](./SPEC.md) (motivation and overview)
- **15 minutes**: add [docs/threat-model.md](./docs/threat-model.md) and [docs/prior-art.md](./docs/prior-art.md)
- **45 minutes**: the full [SPEC.md](./SPEC.md)
- **For NLnet reviewers**: also see [docs/ngi-alignment.md](./docs/ngi-alignment.md)
- **Hands-on (Go)**: `cd go && go test -v ./...` (10 tests pass; see [go/envelope_test.go](./go/envelope_test.go))
- **Hands-on (TypeScript)**: `cd ts && npx tsx --test ./envelope-field.test.ts` (the `TestKnownVectorMatchesGo` test proves the TS output is byte-identical to Go)
- **Working demo**: `cd go && go run ./examples/encrypt_decrypt` (encrypts three envelope fields, decrypts them, demonstrates tamper rejection; see [go/examples/](./go/examples/))

## Comparison with existing systems

| System | Body | Subject | Sender | Recipients | Timestamps | Message-ID |
|---|---|---|---|---|---|---|
| Standard SMTP | Plain | Plain | Plain | Plain | Plain | Plain |
| Proton Mail | Encrypted | Plain | Plain | Plain | Plain | Plain |
| Tuta | Encrypted | Encrypted | Plain | Plain | Plain | Plain |
| **This protocol** | Encrypted | Encrypted | Encrypted | Encrypted | Encrypted | Hashed + encrypted |

(As of July 2026; see SPEC.md Appendix A for sources.)

## License

[Apache 2.0](./LICENSE).

## References

- NIST FIPS 203 (ML-KEM): https://csrc.nist.gov/pubs/fips/203/final
- RFC 7748 (X25519): https://datatracker.ietf.org/doc/html/rfc7748
- RFC 5869 (HKDF): https://datatracker.ietf.org/doc/html/rfc5869
- RFC 5116 (AEAD interface): https://datatracker.ietf.org/doc/html/rfc5116
- RFC 5322 (Internet Message Format): https://datatracker.ietf.org/doc/html/rfc5322
