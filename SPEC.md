# Envelope Metadata Encryption for SMTP-Based Mail Systems

**Version**: 0.1 (draft)
**Status**: Specification phase, pre-implementation
**License**: Apache 2.0

---

## 1. Introduction

### 1.1 Motivation

Modern encrypted email systems have made significant progress in protecting message contents. Both Proton Mail and Tuta encrypt the body, subject, and attachments end-to-end. However, all major encrypted-email providers persist envelope metadata — sender address, recipient addresses, CC/BCC lists, Message-ID, original date, Reply-To, and routing headers — in plaintext at rest.

This metadata is the primary input to traffic analysis. It reveals:

- Who an individual communicates with
- Frequency and timing of communication
- Communication patterns suggestive of relationships, transactions, employment, illness, or activism
- Network graphs of organizations under surveillance

Empirical work by intelligence-agency disclosures and academic researchers establishes that envelope metadata is more valuable to mass-surveillance pipelines than message contents, both because it is structured (machine-tractable) and because it is rarely encrypted.

This specification defines a protocol for encrypting envelope metadata at rest in SMTP-based mail systems, while remaining compatible with the SMTP protocol for external delivery.

### 1.2 Scope

This specification covers:

- Definition of "envelope metadata" for the purposes of this protocol
- Per-email symmetric encryption with authenticated encryption (AES-256-GCM)
- Per-recipient session key wrapping using a hybrid post-quantum key encapsulation mechanism (KEM)
- Wire format for stored encrypted records
- Wire format for the wrapped session key blob
- Cryptographic primitives, parameters, and constants
- Threat model and security considerations

### 1.3 Out of Scope

This specification does NOT cover:

- SMTP protocol implementation or extensions
- Body encryption (assumed delegated; common practice is to encrypt the body with the same session key)
- Attachment encryption (same)
- Key management, key rotation, key recovery
- User authentication
- External (non-Vernam-style) recipient handling
- Storage backend or schema implementation choices
- Indexing, search, or threading semantics for encrypted metadata

### 1.4 Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHOULD", "SHOULD NOT", and "MAY" are to be interpreted as described in RFC 2119.

- **Envelope metadata**: the set of fields defined in §3 of this specification.
- **Session key**: a 32-byte symmetric key generated per email, used to encrypt all envelope fields and (typically) the message body.
- **Wrapped session key**: a session key that has been encapsulated to a specific recipient using the hybrid KEM defined in §4.
- **Recipient**: an entity (typically a user) with a public-key identity, capable of decapsulating wrapped session keys.

---

## 2. Overview

The protocol operates as follows:

1. The sender generates a fresh 32-byte session key using a cryptographically secure random number generator.
2. The sender encrypts each envelope field independently using AES-256-GCM with the session key. Each field receives a fresh 12-byte initialization vector (IV).
3. For each recipient (To/CC/BCC), the sender wraps the session key using the hybrid KEM defined in §4. The recipient's public keys are obtained out of band.
4. The sender stores or transmits:
   - The encrypted envelope fields
   - The set of wrapped session keys (one per recipient)
   - A self-encrypted copy of the session key for the sender's own access (sent items)

5. A recipient decrypts by:
   - Locating their wrapped session key in the set
   - Decapsulating it with their private keys to recover the session key
   - Decrypting each envelope field independently using AES-256-GCM

The server storing the encrypted records observes:
- The number of recipients (length of the wrapped key set)
- Approximate sizes of encrypted fields
- The fact that an email was stored

The server does NOT observe:
- Sender address (encrypted)
- Recipient addresses (encrypted; routing happens via account-id linkage at a higher layer)
- Message-ID (encrypted; SHA-256 hash retained for threading; see §6)
- Subject (encrypted)
- Body (encrypted)
- Original date or received timestamp (encrypted)

---

## 3. Envelope Metadata Fields

The following fields MUST be encrypted under the per-email session key when present:

| Field | Source | RFC | Type |
|---|---|---|---|
| `from` | Header `From:` | RFC 5322 | string |
| `to` | Header `To:` | RFC 5322 | string (comma-separated) |
| `cc` | Header `Cc:` | RFC 5322 | string |
| `bcc` | Envelope RCPT TO | RFC 5321 | string |
| `reply_to` | Header `Reply-To:` | RFC 5322 | string |
| `sender` | Header `Sender:` | RFC 5322 | string |
| `subject` | Header `Subject:` | RFC 5322 | string |
| `message_id` | Header `Message-ID:` | RFC 5322 | string |
| `in_reply_to` | Header `In-Reply-To:` | RFC 5322 | string |
| `references` | Header `References:` | RFC 5322 | string |
| `date` | Header `Date:` | RFC 5322 | RFC 5322 date-time string |
| `received_at` | Server-assigned | implementation | RFC 3339 timestamp |
| `display_name` | parsed from `From` | derived | string |

Implementations MAY encrypt additional fields. Implementations MUST NOT encrypt fields that are required in plaintext for protocol operation (see §6).

---

## 4. Cryptographic Construction

### 4.1 Primitives

The following primitives MUST be used:

| Purpose | Algorithm | Reference |
|---|---|---|
| Authenticated symmetric encryption | AES-256-GCM | NIST SP 800-38D |
| Post-quantum KEM | ML-KEM-1024 | NIST FIPS 203 |
| Classical KEM | X25519 | RFC 7748 |
| Key derivation | HKDF-SHA256 | RFC 5869 |
| Hash | SHA-256 | FIPS 180-4 |
| RNG | Implementation CSPRNG | implementation-defined |

Implementations MUST NOT substitute weaker primitives. ML-KEM-512 and ML-KEM-768 are NOT compliant for this version of the specification. The choice of ML-KEM-1024 is motivated by long-lived email's exposure to "harvest now, decrypt later" attacks; a higher security margin is appropriate.

### 4.2 Per-Email Session Key

For each email, the sender:

1. Generates a 32-byte session key `K_session` from the implementation's CSPRNG.
2. Encrypts each envelope field `F_i` as:

   ```
   IV_i ← random(12)
   ciphertext_i, tag_i ← AES-256-GCM-Encrypt(K_session, IV_i, plaintext_i, AAD = empty)
   stored_i ← IV_i ‖ ciphertext_i ‖ tag_i
   ```

The session key `K_session` MUST NOT be reused across emails.

### 4.3 Hybrid Key Encapsulation

For each recipient R with public-key identity `(KyberPK_R, X25519PK_R)`, the sender wraps the session key as follows:

```
1. (kyber_ct, kyber_ss) ← ML-KEM-1024-Encapsulate(KyberPK_R)
2. (eph_sk, eph_pk)     ← X25519-KeyGen()
3. x25519_ss            ← X25519-DH(eph_sk, X25519PK_R)
4. K_combined           ← HKDF-Expand(
                            ikm  = kyber_ss ‖ x25519_ss,
                            salt = nil,
                            info = "VernamMail-EnvelopeEncryption-v1",
                            L    = 32
                          )
5. iv_wrap              ← random(12)
6. AAD                  ← kyber_ct ‖ eph_pk
7. wrapped, tag         ← AES-256-GCM-Encrypt(K_combined, iv_wrap, K_session, AAD)
8. wire                 ← 0x02 ‖ kyber_ct ‖ eph_pk ‖ iv_wrap ‖ wrapped ‖ tag
```

Output `wire` is exactly **1661 bytes**:

| Bytes | Field |
|---|---|
| 1 | Version (0x02) |
| 1568 | ML-KEM-1024 ciphertext |
| 32 | Ephemeral X25519 public key |
| 12 | AES-GCM IV |
| 32 | Wrapped session key (encrypted) |
| 16 | AES-GCM authentication tag |

### 4.4 Hybrid Decapsulation

Upon receipt, the recipient unwraps the session key as follows:

```
1. Parse wire as: version, kyber_ct, eph_pk, iv_wrap, wrapped, tag.
2. If version ≠ 0x02 (and not v1; see §4.5): reject.
3. kyber_ss     ← ML-KEM-1024-Decapsulate(KyberSK_R, kyber_ct)
4. x25519_ss    ← X25519-DH(X25519SK_R, eph_pk)
5. K_combined   ← HKDF-Expand(
                    ikm  = kyber_ss ‖ x25519_ss,
                    salt = nil,
                    info = "VernamMail-EnvelopeEncryption-v1",
                    L    = 32
                  )
6. AAD          ← kyber_ct ‖ eph_pk
7. K_session    ← AES-256-GCM-Decrypt(K_combined, iv_wrap, wrapped ‖ tag, AAD)
8. Use K_session to decrypt envelope fields per §4.2.
```

If any step fails, the recipient MUST reject the wrapped key and MUST NOT proceed to envelope decryption.

### 4.5 Backwards Compatibility (v1)

A prior version (0x01) of this construction omitted the AAD binding in step 6 of §4.3 and §4.4. Implementations MAY accept v1 wrapped keys for backwards-compatible decryption by setting `AAD = empty` in step 6 of §4.4. New encryptions MUST use v2 (0x02).

The AAD binding (v2) prevents an attacker from substituting `kyber_ct` or `eph_pk` after the wrapping is computed, since the GCM tag would no longer verify.

### 4.5.1 Legacy HKDF Info String

The reference implementation of this protocol originated in a private deployment of Vernam Mail. Prior to the system's 2025 rebrand from "Enigma Inbox" to "Vernam Mail," that deployment used the HKDF info string `"ENIGMA-HybridKEM-v1"`. This canonical specification adopts the clean, vendor-neutral name `"VernamMail-EnvelopeEncryption-v1"` (as specified in §4.3 step 4 and §4.4 step 5).

**For implementations interoperating with the legacy deployment:**

- New wrappings produced by spec-compliant implementations MUST use `"VernamMail-EnvelopeEncryption-v1"`.
- Decryption implementations MAY support a "legacy mode" that, on failure to decrypt with the canonical info string, retries with the legacy info string `"ENIGMA-HybridKEM-v1"`. This fallback is intended solely for migrating ciphertexts produced before the canonical name was adopted.
- The reference Go library exposes legacy-mode decryption as an explicit, off-by-default option. Production deployments are expected to migrate stored ciphertexts to the canonical info string within their normal re-encryption operations (e.g., re-keying, key rotation) and disable legacy mode thereafter.

**Future versions** of this specification (v2 and later) MUST NOT recognize the legacy info string. Migration is REQUIRED before any deployment claims compliance with v1.0 stable of this specification.

This subsection exists to document the historical artifact transparently rather than embed product naming churn into the specification.

---

## 5. Wire Formats

### 5.1 Wrapped Session Key (1661 bytes)

| Offset | Length | Field |
|---|---|---|
| 0 | 1 | Version (0x02) |
| 1 | 1568 | ML-KEM-1024 ciphertext |
| 1569 | 32 | Ephemeral X25519 public key |
| 1601 | 12 | AES-GCM IV |
| 1613 | 32 | Wrapped session key (encrypted) |
| 1645 | 16 | AES-GCM authentication tag |

### 5.2 Encrypted Envelope Field

Each encrypted envelope field is stored as:

| Offset | Length | Field |
|---|---|---|
| 0 | 12 | IV |
| 12 | n | Ciphertext (length = plaintext length) |
| 12 + n | 16 | AES-GCM authentication tag |

Total storage size per field: plaintext length + 28 bytes.

Implementations MAY base64-encode the entire structure for storage in text columns; binary storage is RECOMMENDED.

### 5.3 Stored Email Record (informative)

A typical stored email record contains:

- One encrypted envelope field per metadata column (per §3)
- One wrapped session key per recipient (the recipient's account ID is the lookup key)
- One wrapped session key for the sender (using their own keypair, for sent-items access)
- An encrypted body and attachments (using the same `K_session`; out of scope for this specification)

---

## 6. Operational Plaintext

The following information CANNOT be fully encrypted while maintaining SMTP interoperability and basic mail system function. Implementations MUST be honest about these limitations.

### 6.1 Recipient Account ID

To deliver a stored email to a recipient, the server MUST be able to map the wrapped session key to the recipient's account. Implementations SHOULD use an opaque account ID rather than the recipient's email address for this purpose. The mapping table itself MAY be encrypted to a separate metadata key; see §7.4.

### 6.2 Message-ID Hash for Threading

RFC 5322 requires plaintext Message-ID handling for `In-Reply-To` and `References` traversal. Implementations MAY:

- Store a SHA-256 hash of the Message-ID alongside the encrypted Message-ID, for equality-based threading lookups.
- Optionally also store the encrypted original Message-ID for clients to display the original.

The hash leaks equality but not the original Message-ID value.

### 6.3 External SMTP Egress

Outbound mail to non-Vernam recipients is delivered via standard SMTP, which requires plaintext envelope (`MAIL FROM`, `RCPT TO`) on the wire. Implementations MUST NOT persist this plaintext on disk after the SMTP session closes; envelope plaintext SHOULD reside in process memory only during the active session.

### 6.4 Server-Assigned Routing Metadata

Implementations may need to retain plaintext: account creation timestamps, server identity, retention-policy state. Each implementation MUST document its operational plaintext explicitly.

---

## 7. Security Considerations

### 7.1 Threat Model

The protocol defends against:

- A passive attacker with read access to the storage backend
- A compelled-disclosure adversary (warrant, subpoena) when the operator does not hold recipient private keys
- Future quantum adversaries decrypting present-day captured data ("harvest now, decrypt later")

The protocol does NOT defend against:

- A compromised client device (private keys can be exfiltrated)
- A recipient whose key is leaked or weak
- Active database-write adversaries observing real-time transactions (correlation attacks at the DB layer)
- Traffic analysis of the underlying SMTP transport between providers
- Out-of-band metadata (DNS queries, network endpoints, billing records)

See [docs/threat-model.md](./docs/threat-model.md) for an expanded treatment.

### 7.2 Hybrid Construction Rationale

The hybrid construction (§4.3) requires breaking BOTH ML-KEM-1024 AND X25519 to recover the session key. This provides:

- Defense if ML-KEM-1024 is found to be vulnerable (post-NIST-finalization cryptanalysis)
- Defense if X25519 is broken by a future quantum adversary

The HKDF combination ensures that the resulting key is no weaker than either input alone.

### 7.3 AAD Binding

The v2 AAD construction (`AAD = kyber_ct ‖ eph_pk`) prevents the following attack:

> An attacker who observes a valid wrapped key cannot substitute their own `kyber_ct` or `eph_pk` (e.g., re-wrapping the same `K_session` to a different recipient's key) because the GCM tag would fail.

This binding is REQUIRED for v2 implementations.

### 7.4 Metadata Indirection (Recipient Mapping)

If implementations store recipient account IDs as the lookup key for wrapped session keys, those account IDs are visible to the server. Implementations MAY add an additional layer:

- The wrapped key set is keyed by HMAC(server_secret, account_id), so the server cannot enumerate all emails for a given account without the secret
- Or recipient mapping is stored encrypted under a metadata key derived from each user's authentication

This is implementation-defined and outside this specification.

### 7.5 Constant-Time Operations

Implementations MUST use constant-time operations for:

- AES-256-GCM authentication tag verification
- Comparison of decapsulated values
- Key material handling

Implementations SHOULD zeroize key material after use where the platform supports it.

### 7.6 Random Number Generation

All randomness (session keys, IVs, ephemeral keys) MUST come from a cryptographically secure RNG. On supported platforms:

- POSIX: `/dev/urandom` or `getrandom(2)`
- Windows: `BCryptGenRandom`
- Web: `crypto.getRandomValues`

### 7.7 Replay and Substitution

Each session key is generated fresh per email. The wrapped key wire format is bound to the specific (`kyber_ct`, `eph_pk`) pair via AAD. There is no defined replay defense at this layer; implementations needing replay protection SHOULD use higher-layer mechanisms (e.g., a signed canonical email hash).

### 7.8 HKDF Salt

The HKDF salt is fixed as nil, with the strength of the derivation resting on the entropy of the IKM (a 1568-byte ML-KEM-1024 shared secret concatenated with a 32-byte X25519 shared secret). Implementations MUST NOT vary the salt.

---

## 8. Compliance and Test Vectors

### 8.1 Compliance

An implementation is compliant with this specification if it:

1. Correctly produces and consumes the wire formats defined in §5
2. Uses only the primitives and parameters defined in §4.1
3. Uses the exact canonical info string `"VernamMail-EnvelopeEncryption-v1"` in HKDF-Expand for all NEW encryptions (legacy-mode decryption per §4.5.1 is permitted but does not constitute spec compliance)
4. Implements the AAD binding defined in §4.3 step 6 for v2 wrapping
5. Documents its operational plaintext per §6.4
6. If supporting legacy-mode decryption (§4.5.1), exposes it as an explicit, off-by-default option and documents migration intent

### 8.2 Test Vectors

See the [`test-vectors/`](./test-vectors/) directory for example inputs and expected outputs. Test vectors will be expanded in v0.2.

---

## 9. Implementation Notes

### 9.1 Library Boundaries

A reference implementation provides:

- Session key generation
- Per-field envelope encryption / decryption with AES-256-GCM
- Hybrid encapsulation / decapsulation
- Wire format encoders / decoders
- Test vector validation

A reference implementation does NOT provide:

- SMTP integration (MTA layer)
- Database schema or persistence
- Key management (storage of recipient public keys)
- User authentication
- Body or attachment encryption (delegated; uses the same session key with the same primitives)

### 9.2 Performance (informative)

The following are approximate, pre-benchmarking estimates for a modern x86_64 CPU, single-threaded. Formal measurements will be published as part of milestone 4 of the reference implementation. Implementations and consumers SHOULD NOT rely on these numbers for capacity planning.

- ML-KEM-1024 encapsulation: ~50 µs
- X25519 ephemeral key generation + DH: ~80 µs
- HKDF-SHA256 expand (32 bytes): ~10 µs
- AES-256-GCM encrypt (1 KB envelope field): ~5 µs
- Estimated total wrap operation: ~150 µs

Estimated per-email overhead for an email with 5 recipients and 12 envelope fields: ~750 µs encryption, ~750 µs decryption per recipient.

### 9.3 Storage Overhead

Per email, the AEAD overhead and wrapped-key footprint sum as follows:

- 12 envelope fields × 28 bytes AEAD overhead (12-byte IV + 16-byte tag) = 336 bytes
- Each wrapped session key: 1661 bytes (one per recipient + one for the sender's own copy)
- For an email with 5 external recipients, total wrapped keys = 6 × 1661 = 9,966 bytes
- Combined per-email overhead (5 external recipients): 336 + 9,966 ≈ **10.3 KB**

For a single-recipient email, the per-email overhead is ~336 + 2 × 1661 ≈ **3.7 KB**. Overhead grows linearly with recipient count at 1661 bytes per recipient. This is small relative to typical message body sizes, where bodies and attachments commonly run from tens to hundreds of kilobytes.

---

## 10. References

### Normative

- [FIPS-203] NIST FIPS 203, "Module-Lattice-Based Key-Encapsulation Mechanism Standard," August 2024
- [RFC-7748] Langley, A., et al., "Elliptic Curves for Security," IETF RFC 7748, January 2016
- [RFC-5869] Krawczyk, H., Eronen, P., "HMAC-based Extract-and-Expand Key Derivation Function (HKDF)," IETF RFC 5869, May 2010
- [RFC-5116] McGrew, D., "An Interface and Algorithms for Authenticated Encryption," IETF RFC 5116, January 2008
- [SP-800-38D] NIST SP 800-38D, "Recommendation for Block Cipher Modes of Operation: Galois/Counter Mode (GCM) and GMAC," November 2007
- [FIPS-180-4] NIST FIPS 180-4, "Secure Hash Standard (SHS)," August 2015

### Informative

- [RFC-5322] Resnick, P., "Internet Message Format," IETF RFC 5322, October 2008
- [RFC-5321] Klensin, J., "Simple Mail Transfer Protocol," IETF RFC 5321, October 2008
- [RFC-2119] Bradner, S., "Key words for use in RFCs to Indicate Requirement Levels," IETF RFC 2119, March 1997
- [ETSI-TR-103-619] ETSI TR 103 619 V1.1.1, "CYBER; Migration strategies and recommendations to Quantum Safe schemes," European Telecommunications Standards Institute, July 2020
- [ETSI-TR-103-823] ETSI TR 103 823 V1.1.1, "CYBER; Quantum-Safe Public-Key Encryption and Key Encapsulation," European Telecommunications Standards Institute, March 2022
- [I-D-IRTF-CFRG-HPKE] Barnes, R., Bhargavan, K., Lipp, B., Wood, C. A., "Hybrid Public Key Encryption," RFC 9180, February 2022
- [STEBILA-HYBRID] Stebila, D., Fluhrer, S., Gueron, S., "Hybrid key exchange in TLS 1.3," IETF draft `draft-ietf-tls-hybrid-design`

---

## Appendix A: Comparison with Existing Systems

| System | Body | Subject | Sender | Recipients | Timestamps | Message-ID |
|---|---|---|---|---|---|---|
| Standard SMTP | Plain | Plain | Plain | Plain | Plain | Plain |
| Proton Mail | Encrypted | Plain | Plain | Plain | Plain | Plain |
| Tuta | Encrypted | Encrypted | Plain | Plain | Plain | Plain |
| Vernam Mail (this spec) | Encrypted | Encrypted | Encrypted | Encrypted | Encrypted | Hashed + encrypted |

(As of April 2026; verify before citing in derivative work.)

---

*End of specification.*
