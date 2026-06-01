# Threat Model

This document expands on §7 of the main specification with explicit attacker capabilities and protocol guarantees.

## Adversaries

### A1: Passive Storage Adversary

**Capabilities:**
- Read-only access to the storage backend (database, files, backups)
- Can dump tables, read encrypted blobs
- Cannot modify data or interact with running services

**Defended against:** YES. All envelope metadata is AES-256-GCM encrypted; A1 sees only ciphertext.

### A2: Active Database Adversary

**Capabilities:**
- Read and write access to the storage backend during operation
- Can modify ciphertexts, replace records, observe transaction-level operations

**Partially defended against:** AES-GCM authentication prevents undetected ciphertext modification. However, A2 can observe write transactions in real time and infer correlations (e.g., "row written to sender's sent table at time T, row written to recipient's inbox at time T+δ"). This is documented as out of scope; see §7.1 of the spec.

### A3: Compelled Disclosure Adversary

**Capabilities:**
- Legal authority to compel the operator to produce stored data
- Subpoena, warrant, NSL, etc.

**Defended against:** PARTIALLY. The operator can produce only encrypted blobs and the (potentially hashed) recipient mapping. The operator cannot produce plaintext envelope metadata for stored emails because they do not possess recipient private keys.

**NOT defended against:** Compelled disclosure of newly-arriving emails before encryption (i.e., during SMTP ingestion); compelled provision of malicious client code.

### A4: Future Quantum Adversary (Harvest Now, Decrypt Later)

**Capabilities:**
- Captures all current encrypted data
- Possesses a sufficiently large quantum computer at some future date
- Attempts to break X25519 via Shor's algorithm

**Defended against:** YES. The hybrid construction requires breaking BOTH ML-KEM-1024 AND X25519. ML-KEM-1024 is believed to remain secure against quantum attack at the AES-256-equivalent level.

### A5: Compromised Client Device

**Capabilities:**
- Full access to the user's private keys via filesystem, memory, or runtime injection

**NOT defended against.** Out of scope. Protecting against a compromised endpoint requires hardware-rooted trust (TEE, HSM) which is outside the protocol layer.

### A6: Compromised Recipient

**Capabilities:**
- The recipient is the adversary; they receive the email and possess valid keys to decrypt it

**NOT defended against.** This is fundamental: the recipient is the intended audience.

### A7: Traffic Analysis Adversary

**Capabilities:**
- Network-level observation of TLS-encrypted SMTP traffic between mail servers
- DNS query observation
- Timing analysis

**NOT defended against.** This protocol operates at the storage layer, not the network layer. Defenses against traffic analysis require lower-layer mechanisms (Tor, mix networks, batching).

## Out-of-Scope Concerns

- Forward secrecy across email retention periods
- Post-compromise security after key rotation
- Verification of sender identity (delegated to signing layers)
- Spam, phishing, abuse mitigation
- Quota and storage management

## Trust Boundaries

The protocol places trust as follows:

| Component | Trusted For |
|---|---|
| Sender's client | Generating session keys; encrypting honestly |
| Recipient's client | Decrypting and rendering honestly |
| Recipient's private keys | Confidentiality; not stolen |
| Server / operator | Storing encrypted blobs honestly; not providing malicious client code |
| TLS layer | Confidentiality of transport (orthogonal to this layer) |

The protocol does NOT trust:

- The server / operator with the contents of the stored data
- The network with anything beyond what TLS provides
- Other recipients of the same email (they can decrypt that email but not other emails)

## Known Limitations

1. **Database-layer transaction correlation.** As noted under A2, an adversary with live write access can correlate sender and recipient inserts within a transaction window. This is a known limitation; mitigations require batching, delays, or external write randomization, all of which are out of scope for this specification.

2. **Operational plaintext.** As noted in §6 of the spec, certain fields cannot be encrypted while maintaining SMTP interoperability and routing function. Implementations MUST disclose their operational plaintext.

3. **Recipient key compromise.** If a recipient's private key is compromised, all past emails encrypted to that recipient become decryptable. This is a fundamental property of public-key encryption; forward secrecy at the email level requires more elaborate constructions.

4. **No verification of sender.** This protocol does not authenticate the sender. An adversary in possession of a recipient's public keys can produce a valid-looking encrypted email from any claimed sender. Implementations SHOULD layer signatures (e.g., DKIM, OpenPGP signatures) for sender authentication.
