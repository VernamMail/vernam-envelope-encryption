# Prior Art and Related Work

This document surveys related work in encrypted email and post-quantum cryptography and explains how this specification relates to it.

## Encrypted Email Systems

### Proton Mail

Proton Mail provides end-to-end encryption for the message body and attachments using OpenPGP (the subject line is not encrypted). Proton has announced work toward post-quantum support via hybrid constructions for the content-encryption layer; see their blog for current rollout status.

**Envelope metadata handling:** Proton Mail stores `From`, `To`, `Cc`, `Bcc`, `Reply-To`, `Message-ID`, and timestamps in plaintext at rest. Search index data (with the user's consent) is also indexed plaintext on the server in some configurations.

**Reference:** https://proton.me/blog/post-quantum-encryption-mail

### Tuta (formerly Tutanota)

Tuta encrypts the body, subject, and attachments. It deployed TutaCrypt, a hybrid post-quantum KEM, in March 2024.

**Envelope metadata handling:** Subject is encrypted (notable; ahead of Proton). Sender and recipient addresses, Message-ID, and timestamps are stored in plaintext.

**Reference:** https://tuta.com/blog/post-quantum-cryptography

### Skiff (decommissioned)

Notion acquired Skiff in 2024 and discontinued the service. Skiff's design was similar to Proton's; envelope metadata was plaintext at rest.

## Post-Quantum Cryptography Standards

### NIST PQC Project

The NIST Post-Quantum Cryptography Standardization Project finalized three key standards in August 2024:

- **FIPS 203**: Module-Lattice-Based Key-Encapsulation Mechanism Standard (ML-KEM, formerly Kyber)
- **FIPS 204**: Module-Lattice-Based Digital Signature Standard (ML-DSA, formerly Dilithium)
- **FIPS 205**: Stateless Hash-Based Digital Signature Standard (SLH-DSA, formerly SPHINCS+)

This specification uses ML-KEM-1024 (FIPS 203) for the post-quantum component of its hybrid KEM.

**Reference:** https://csrc.nist.gov/Projects/post-quantum-cryptography

### Hybrid Construction Patterns

The pattern of combining a post-quantum KEM with a classical KEM (e.g., X25519) and deriving a single shared key via HKDF is well-established. See:

- Stebila, D., et al., "Hybrid key exchange in TLS 1.3," IETF draft `draft-ietf-tls-hybrid-design`
- ETSI TR 103 619, "Migration strategies and recommendations to Quantum Safe schemes"
- Barnes, R., Bhargavan, K., Lipp, B., Wood, C. A., "Hybrid Public Key Encryption," RFC 9180, February 2022

### NLnet-Funded PQC Projects (Adjacent)

NLnet has funded several post-quantum cryptography projects relevant to this work:

- **Quantum-Safe Cryptography in Sequoia PGP** (NGI0 Commons Fund, 2025): implements `draft-ietf-openpgp-pqc` in Sequoia PGP
- **CurveForge** (NGI0 Commons Fund, 2026): optimized post-quantum arithmetic toolkit
- **Rosenpass** (NGI Assure, 2022-2024): post-quantum security add-on for WireGuard
- **oqsprovider** (NGI Assure, 2021-2023): post-quantum algorithms for OpenSSL
- **KEMTLS standardization** (NGI Assure, 2021-2023)

This specification aims to complement these efforts by addressing the email-system layer specifically, with emphasis on envelope metadata rather than message body.

## What This Specification Adds

To the best of our knowledge, the combination of:

1. **Per-field envelope metadata encryption at rest** (encrypting From, To, CC, BCC, Message-ID, timestamps independently with AES-256-GCM)
2. **Per-recipient session key wrapping via hybrid post-quantum KEM** with ciphertext binding
3. **Compatibility with standard SMTP for external delivery**

…is not addressed by any prior published specification. Existing encrypted-email systems either:

- Encrypt only message contents (Proton, Tuta)
- Operate as walled gardens incompatible with SMTP (Bitmessage, certain decentralized systems)
- Use signatures/encryption at the transport layer (DKIM, S/MIME) which do not encrypt at rest

This specification is intended to fill that gap as an open commons contribution that other privacy-preserving email systems can adopt.

## OpenPGP Considerations

This specification does not replace OpenPGP. OpenPGP and S/MIME remain valuable for end-to-end encryption between systems that share the same standard. This specification operates at a different layer: encrypting metadata at rest within a single mail system's own storage.

Implementations MAY interoperate with OpenPGP for cross-provider encrypted mail by:

- Using OpenPGP for message body encryption when communicating with providers that do not implement this protocol
- Re-wrapping inbound OpenPGP-encrypted messages under this protocol's session-key model for at-rest storage
- Maintaining OpenPGP keys solely for external interoperability, never for internal storage encryption

## Academic References (selected)

- Schneier, B. "Data and Goliath: The Hidden Battles to Collect Your Data and Control Your World," W.W. Norton, 2015. On the value of metadata to surveillance.
- Greenwald, G. "No Place to Hide," Metropolitan Books, 2014. On NSA bulk metadata collection.
- Bernstein, D., Lange, T. "Post-quantum cryptography," Nature, 2017
- "Mass Surveillance," EFF Surveillance Self-Defense Project: https://ssd.eff.org/

## Email Privacy Coalition (informative)

Emerging coalitions and standards bodies working on related problems:

- IETF MAILMAINT working group
- IRTF CFRG (Crypto Forum Research Group)
- IETF LAMPS (Limited Additional Mechanisms for PKIX and SMIME)

This specification's authors intend to engage these communities as part of the spec's evolution.
