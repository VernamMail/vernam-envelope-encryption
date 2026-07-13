# TypeScript reference: envelope-field encryption

A minimal, standalone TypeScript reference implementation of the envelope-**field** encryption portion of the protocol (SPEC.md §4.2 and §5.2).

It depends only on the Web Crypto API (`globalThis.crypto.subtle`), with no framework, storage, or product dependencies, so it runs unchanged in browsers, Node.js (≥ 20), Deno, and Bun.

## Why this exists

The protocol runs in Vernam Mail's pre-launch deployment: the Go server and TypeScript web client carry all pre-launch traffic, and a Kotlin (Android) implementation interoperates in ongoing integration testing. This module is the public, vendor-neutral TypeScript reference for the envelope-field layer. It exists primarily to **demonstrate that the specification is language-agnostic**: it produces byte-identical output to the [Go reference](../go/envelope.go) for the same inputs.

The cross-language test (`TestKnownVectorMatchesGo` in `envelope-field.test.ts`) asserts that this implementation reproduces test vector `envelope-field-001` from [../test-vectors/basic.json](../test-vectors/basic.json), the same vector verified by the Go test suite. Same key, same IV, same plaintext, identical wire bytes.

## Scope

This module implements **only** envelope-field encryption:

- `newSessionKey()`: generate a 32-byte session key
- `encryptField(key, plaintext)`: AES-256-GCM with a fresh random IV
- `decryptField(key, ciphertext)`: authenticated decryption
- `encryptFieldWithIV(key, iv, plaintext)`: deterministic core (testing only)

It does **not** implement the hybrid ML-KEM-1024 + X25519 key wrapping specified in SPEC.md §4.3 (`Wrap`/`Unwrap` are declared but stubbed in the Go reference as well). That is scheduled as a cross-language deliverable of the roadmap; see [../ROADMAP.md](../ROADMAP.md) milestone 3.

## Run the tests

```sh
cd ts
npm install      # installs tsx + typescript (dev only)
npm test         # runs the cross-language vector match + round-trip tests
```

Or directly, without installing:

```sh
npx tsx --test ./envelope-field.test.ts
```

## Wire format

Per SPEC.md §5.2, each encrypted field is:

```
IV (12 bytes) || ciphertext (n bytes) || GCM tag (16 bytes)
```

This matches the Go reference exactly. Total size = plaintext length + 28 bytes.

## License

[Apache 2.0](../LICENSE).
