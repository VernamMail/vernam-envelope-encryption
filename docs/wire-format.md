# Wire Formats

This document gives a concrete byte-level layout for the wire formats defined in §5 of the specification, with worked diagrams.

## Wrapped Session Key (1661 bytes total)

```
 0                                                        7
 +---+----------------+---------------+----+--------------+----+
 | V |   kyber_ct     |    eph_pk     | IV |  wrapped_sk  |TAG |
 +---+----------------+---------------+----+--------------+----+
   1        1568            32          12       32         16

 V        = version byte (0x02 for current)
 kyber_ct = ML-KEM-1024 ciphertext
 eph_pk   = ephemeral X25519 public key
 IV       = 12-byte AES-GCM nonce
 wrapped_sk = AES-256-GCM-encrypted session key (32 bytes)
 TAG      = 16-byte AES-GCM authentication tag
```

### Byte Offsets

| Offset (decimal) | Offset (hex) | Length | Field |
|---|---|---|---|
| 0 | 0x000 | 1 | Version |
| 1 | 0x001 | 1568 | ML-KEM-1024 ciphertext |
| 1569 | 0x621 | 32 | Ephemeral X25519 public key |
| 1601 | 0x641 | 12 | AES-GCM IV |
| 1613 | 0x64D | 32 | Wrapped session key |
| 1645 | 0x66D | 16 | AES-GCM authentication tag |
| **1661** | **0x67D** | - | (end) |

## Encrypted Envelope Field (variable length)

```
 0                            n+12        n+28
 +----+---------------------+----+
 | IV |     ciphertext      |TAG |
 +----+---------------------+----+
  12          n               16

 IV         = 12-byte AES-GCM nonce
 ciphertext = AES-256-GCM-encrypted plaintext, length = n bytes
 TAG        = 16-byte AES-GCM authentication tag
```

Total length: `n + 28` bytes where `n` is plaintext length.

### Encoding

Implementations MAY base64-encode the entire structure for storage in text columns. RECOMMENDED: store as raw bytes (e.g., `BLOB` / `BYTEA`).

## Worked Example (informative)

For an email with the From header `alice@example.com` (17 bytes) and one recipient:

- Encrypted `from` field: 17 + 28 = 45 bytes
- Wrapped session key (1 recipient): 1661 bytes
- Wrapped session key for sender: 1661 bytes (so sender can decrypt sent items)

For all 12 envelope fields with similar sizes (~30 bytes plaintext average):

- Total encrypted fields: ~12 × (30 + 28) = ~696 bytes
- Wrapped keys (1 recipient + sender): 2 × 1661 = ~3.3 KB
- **Total per-email envelope overhead: ~4 KB**

For an email with 10 recipients:

- Encrypted fields: ~696 bytes
- Wrapped keys: 11 × 1661 = ~18.3 KB
- **Total: ~19 KB**

## Endianness and Encoding

All multi-byte integer fields (none are present in v0.1, but reserved for future extensions) MUST use big-endian (network) byte order.

All key material is treated as opaque byte sequences. ML-KEM-1024 ciphertexts and X25519 public keys are stored exactly as their respective specifications produce them.
