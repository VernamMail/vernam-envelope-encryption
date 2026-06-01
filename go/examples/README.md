# Examples

Working programs demonstrating the envelope-encryption library.

| Directory | Description |
|---|---|
| `encrypt_decrypt/` | Minimal end-to-end demo of envelope-field encryption, decryption, and AEAD tamper rejection. |

## Running

From the Go module root (one level up from this directory):

```sh
go run ./examples/encrypt_decrypt
```

Additional examples covering hybrid KEM wrap/unwrap and per-recipient session-key distribution will be added during NLnet milestones 3 and 4.
