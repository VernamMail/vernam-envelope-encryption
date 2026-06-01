# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in this specification or in the reference implementation, please report it privately. Do **not** open a public GitHub issue.

**Preferred channel:** email `security@vernammail.com` with the subject prefix `[envelope-encryption]`.

If your report contains sensitive information (proof-of-concept exploit, key material, etc.), please request our PGP public key in a clear-text email and we will respond with the key and fingerprint via a separate channel. A long-term published PGP key location will be added here when available.

We aim to acknowledge receipt within **72 hours** and provide an initial assessment within **7 days**. Resolution timelines depend on severity and complexity.

## Scope

In scope:

- Cryptographic flaws in the protocol specification (SPEC.md), including but not limited to:
  - Weaknesses in the hybrid KEM construction (§4.3)
  - AAD-binding bypass scenarios (§7.3)
  - Wire-format ambiguities that allow ciphertext substitution
  - Side-channel leaks identifiable from the spec
- Cryptographic flaws in the reference implementation (`go/`):
  - Incorrect primitive usage (e.g., nonce reuse, missing authentication)
  - Timing or oracle leaks
  - Memory-safety issues affecting key material handling
- Documentation errors that would lead an honest implementer to produce an insecure system

Out of scope:

- Vulnerabilities in dependencies (Go standard library, `crypto/mlkem`, `crypto/ecdh`) — please report those upstream
- Attacks already documented as out-of-scope in [docs/threat-model.md](./docs/threat-model.md), including:
  - Compromised client devices
  - Active database-write adversaries causing transaction-level correlation
  - Traffic analysis at the SMTP transport layer
- Issues in third-party software that integrates this protocol — report those to the integrator

## Coordinated Disclosure

We follow a 90-day coordinated disclosure timeline by default:

1. **Day 0:** report received, acknowledged within 72 hours
2. **Day 1–14:** assessment, root-cause analysis, fix design
3. **Day 14–60:** patch development, testing, security advisory drafted
4. **Day 60–90:** patch released to users, public advisory published

Earlier public disclosure may be agreed by mutual consent if a fix is straightforward. Later disclosure may be requested if the fix is complex or coordinated with other affected projects (e.g., upstream cryptographic libraries).

## What to Expect

- We will not threaten legal action against good-faith security researchers
- We will credit reporters in the published advisory unless they prefer anonymity
- We will not pay bug bounties at this time (pre-launch, no funding for that program)
- We will not consider the following as security issues:
  - Theoretical attacks requiring resources clearly outside the threat model
  - Issues in code paths explicitly marked as out of scope (e.g., demos, ignored build files)
  - Best-practice suggestions without a concrete attack scenario (please open a regular issue or PR)

## Pre-Implementation Phase

This repository is currently in specification phase. The reference implementation is partial. The following types of finding are particularly welcome:

- Errors or ambiguities in SPEC.md
- Missing edge cases in the threat model
- Concerns about the choice of primitives or parameters
- Weaknesses in the AAD-binding construction (v2)

You may also contact us with general cryptographic-design feedback that is not strictly a vulnerability; that traffic can go to a regular issue or pull request.

---

Last updated: 2026-05-02
