# NGI Alignment and European Dimension

This document maps the envelope-metadata-encryption protocol's technical contribution to the goals of the European Commission's [Next Generation Internet (NGI) initiative](https://www.ngi.eu/) and addresses the "European Dimension" eligibility criterion of the NLnet Restack fund for non-EU applicants.

## Project Context

The applicant team consists of a Canadian-incorporated solo founder (Imporon Inc., operating as Vernam Mail). The proposal seeks Restack fund support to extract a privacy-preserving email-encryption protocol from a private codebase and publish it as an open commons artifact: specification, reference Go library, cross-language interop test vectors, and third-party cryptographic review. Research and development is the primary objective: the deliverables are a formal protocol specification, a reference implementation, an interoperability test suite validated across three independent implementations, and an independent cryptographic review, not a commercial product.

Although the applicant is non-EU, the project's substantive contribution falls within the NGI program's stated mission, with several concrete European-aligned characteristics enumerated below.

## NGI Mission Alignment

The NGI initiative's stated mission is to "shape the development of the Internet of tomorrow as an Internet that responds to people's fundamental needs, including trust, security, and inclusion." This project advances each of those pillars directly.

### Trust

Existing encrypted email systems (Proton Mail, Tuta) build user trust through end-to-end encryption of message contents. They simultaneously erode that trust by persisting envelope metadata in plaintext at rest, where it can be read by operators or compelled-disclosure adversaries. This protocol closes that gap, allowing operators to make a stronger, more honest claim: *"We cannot see who you email, when, or how often."*

### Security

The hybrid post-quantum construction (ML-KEM-1024 + X25519 with HKDF combination and ciphertext-binding AAD) provides defense against both classical and future quantum adversaries. The "harvest now, decrypt later" threat is particularly acute for email, where messages may retain sensitivity for decades. ML-KEM-1024 is NIST FIPS 203-standardized and corresponds in security level to AES-256 against quantum cryptanalysis. The hybrid design ensures continued security even if either constituent primitive is later weakened.

### Inclusion

The protocol is published under Apache 2.0 with a vendor-neutral specification. Any encrypted-email provider, federated mail system, or research project may adopt it without licensing barriers. The specification is intentionally protocol-agnostic with respect to storage backend, key management, and authentication, allowing diverse implementers, from large providers to self-hosted single-user systems, to incorporate it.

## European Standards and Research

The specification cites and aligns with European cryptographic standards alongside NIST and IETF standards:

- **ETSI TR 103 619**, "Migration strategies and recommendations to Quantum Safe schemes": the ETSI guidance on transitioning to post-quantum cryptography. This protocol's hybrid construction is consistent with the migration patterns recommended therein.
- **ETSI TR 103 823**, "Quantum-Safe Public-Key Encryption and Key Encapsulation": directly applicable to the hybrid KEM construction.
- **IETF RFC 9180 (HPKE)**, Hybrid Public Key Encryption: informs the structure of our combination of post-quantum and classical KEMs.
- **Stebila et al., "Hybrid key exchange in TLS 1.3"** (`draft-ietf-tls-hybrid-design`): design pattern reference.

Researchers at European institutions whose foundational work this protocol builds upon include Daniel J. Bernstein (TU Eindhoven, designer of Curve25519/X25519), Peter Schwabe (Radboud University and MPI-SP, co-author of CRYSTALS-Kyber, standardized as ML-KEM), and the KU Leuven COSIC group (designers of the SABER NIST post-quantum finalist and a long-standing center of lattice-cryptography research).

## Privacy-by-Design and GDPR Alignment

The protocol is constructed to minimize the data an operator can be compelled to disclose, by ensuring that the operator does not hold the keys necessary to decrypt envelope metadata for stored emails. This aligns with several principles of EU GDPR and the broader European data-protection tradition:

- **Data minimization** (GDPR Art. 5(1)(c)): the operator cannot retain plaintext metadata it does not require for routing
- **Purpose limitation** (GDPR Art. 5(1)(b)): metadata not stored as plaintext cannot be repurposed for analytics or profiling
- **Storage limitation** (GDPR Art. 5(1)(e)): encrypted metadata at rest reduces the footprint of identifiable data even when retention is long
- **Security of processing** (GDPR Art. 32): the protocol uses AES-256-GCM and a hybrid post-quantum KEM to encrypt metadata at rest

The protocol is also consistent with the direction of the EU's Cyber Resilience Act, which emphasizes secure defaults, transparency, and resilience against future cryptographic threats.

## Hosting and Deployment Geography

The reference deployment (Vernam Mail's pre-launch system) is hosted on infrastructure in France (OVH), within EU jurisdiction. Production data resides in the EU under GDPR. While this is a deployment characteristic of the existing private system rather than a property of the open protocol itself, it demonstrates that the protocol has been built and tested in a real European-hosted environment.

## Beneficiaries

The open commons deliverable benefits every implementer of encrypted email globally, including European users, providers, and researchers. Specifically:

- **European privacy-preserving email providers** (Tuta, Mailfence, Mailbox.org, Disroot, RiseUp, etc.) gain a vendor-neutral specification they may adopt to reduce metadata exposure beyond their current designs.
- **European researchers** working on email privacy and post-quantum cryptography gain a documented protocol and reference implementation suitable for academic citation, extension, and analysis.
- **European NGOs and journalists** whose threat model includes metadata-based surveillance benefit from any implementer adopting the protocol.
- **The European Internet community** gains an open contribution to the privacy-preserving infrastructure layer that NGI is explicitly designed to advance.

## Coordination with NLnet-Funded Work

This project is designed to complement, not duplicate, prior and concurrent NLnet-funded post-quantum work:

- **Quantum-Safe Cryptography in Sequoia PGP** (NGI0 Commons Fund, 2025): operates at the OpenPGP layer for message-content encryption. This project operates at the envelope-metadata layer in mail-system storage, a distinct concern.
- **Rosenpass** (NGI Assure, 2022-2024): post-quantum extension to WireGuard, transport-layer. This project operates at the application-layer storage encryption, orthogonal.
- **oqsprovider** (NGI Assure, 2021-2023): post-quantum primitives for OpenSSL. This project consumes such primitives but does not itself provide them.
- **CurveForge** (NGI0 Commons Fund, 2026): low-level post-quantum arithmetic. Adjacent; could be used at the implementation layer.

The proposed deliverable fills a gap not addressed by these existing projects: encryption of email envelope metadata at the storage layer, with a vendor-neutral specification and reference implementation.

## Summary

The European Dimension here is in the work, not the applicant's address. The deliverable is an open, vendor-neutral protocol that European email providers, researchers, and civil-society organizations can adopt directly. The applicant happens to be incorporated in Canada; the commons it produces is available to everyone, EU included.
