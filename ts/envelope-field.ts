// Package envelope-field provides the TypeScript reference implementation of
// the envelope-FIELD encryption portion of the protocol described in
// ../SPEC.md (sections 4.2 and 5.2).
//
// This module is intentionally minimal and standalone: it depends only on the
// Web Crypto API (globalThis.crypto.subtle), with no framework, storage, or
// product dependencies. It mirrors the Go reference in ../go/envelope.go so
// that both implementations produce byte-identical output for the same inputs
// (verified against test-vectors/basic.json, vector envelope-field-001).
//
// The hybrid ML-KEM-1024 + X25519 key wrapping (Wrap/Unwrap in the Go
// reference) is NOT included here; it is scheduled as a funded cross-language
// deliverable. See ../ROADMAP.md milestone 3.

// Protocol constants. See SPEC.md sections 4.2 and 5.2.
export const SESSION_KEY_SIZE = 32; // bytes (AES-256)
export const GCM_NONCE_SIZE = 12; // bytes
export const GCM_TAG_SIZE = 16; // bytes

/** A 32-byte symmetric session key, generated per email. */
export type SessionKey = Uint8Array;

/**
 * Wire format for a single encrypted envelope field, per SPEC.md section 5.2:
 *   IV (12 bytes) || ciphertext (n bytes) || tag (16 bytes)
 */
export type EncryptedField = Uint8Array;

export class MalformedWireError extends Error {
  constructor() {
    super("envelope: malformed wire format");
    this.name = "MalformedWireError";
  }
}

export class AuthenticationError extends Error {
  constructor() {
    super("envelope: authentication failed");
    this.name = "AuthenticationError";
  }
}

export class InvalidNonceSizeError extends Error {
  constructor() {
    super("envelope: invalid nonce size");
    this.name = "InvalidNonceSizeError";
  }
}

/** Generate a fresh 32-byte session key from the platform CSPRNG. */
export function newSessionKey(): SessionKey {
  return crypto.getRandomValues(new Uint8Array(SESSION_KEY_SIZE));
}

async function importKey(
  key: SessionKey,
  usage: KeyUsage,
): Promise<CryptoKey> {
  return crypto.subtle.importKey("raw", key, { name: "AES-GCM" }, false, [
    usage,
  ]);
}

/**
 * Encrypt a single envelope field under the given session key using
 * AES-256-GCM with a fresh random IV. Output layout per SPEC.md section 5.2:
 * IV || ciphertext || tag.
 */
export async function encryptField(
  key: SessionKey,
  plaintext: Uint8Array,
): Promise<EncryptedField> {
  const iv = crypto.getRandomValues(new Uint8Array(GCM_NONCE_SIZE));
  return encryptFieldWithIV(key, iv, plaintext);
}

/**
 * Deterministic core used by encryptField and by test-vector generators.
 * Exposed for testing only: production callers MUST use encryptField, which
 * supplies a fresh random IV. Reusing an IV under the same key breaks
 * AES-GCM's confidentiality and integrity.
 */
export async function encryptFieldWithIV(
  key: SessionKey,
  iv: Uint8Array,
  plaintext: Uint8Array,
): Promise<EncryptedField> {
  if (iv.length !== GCM_NONCE_SIZE) {
    throw new InvalidNonceSizeError();
  }
  const cryptoKey = await importKey(key, "encrypt");
  const sealed = new Uint8Array(
    await crypto.subtle.encrypt({ name: "AES-GCM", iv }, cryptoKey, plaintext),
  );
  // Web Crypto returns ciphertext || tag concatenated, matching Go's Seal.
  const out = new Uint8Array(iv.length + sealed.length);
  out.set(iv, 0);
  out.set(sealed, iv.length);
  return out;
}

/**
 * Decrypt a previously encrypted envelope field.
 * Throws MalformedWireError if the input is too short to contain IV + tag,
 * or AuthenticationError if the GCM tag does not verify.
 */
export async function decryptField(
  key: SessionKey,
  ciphertext: EncryptedField,
): Promise<Uint8Array> {
  if (ciphertext.length < GCM_NONCE_SIZE + GCM_TAG_SIZE) {
    throw new MalformedWireError();
  }
  const iv = ciphertext.subarray(0, GCM_NONCE_SIZE);
  const ct = ciphertext.subarray(GCM_NONCE_SIZE);
  const cryptoKey = await importKey(key, "decrypt");
  try {
    const plaintext = await crypto.subtle.decrypt(
      { name: "AES-GCM", iv },
      cryptoKey,
      ct,
    );
    return new Uint8Array(plaintext);
  } catch {
    throw new AuthenticationError();
  }
}

/** Hex-encode bytes (lowercase). Utility for test vectors and debugging. */
export function toHex(bytes: Uint8Array): string {
  return Array.from(bytes)
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("");
}

/** Decode a lowercase/uppercase hex string to bytes. */
export function fromHex(hex: string): Uint8Array {
  if (hex.length % 2 !== 0) {
    throw new Error("envelope: odd-length hex string");
  }
  const out = new Uint8Array(hex.length / 2);
  for (let i = 0; i < out.length; i++) {
    out[i] = parseInt(hex.slice(i * 2, i * 2 + 2), 16);
  }
  return out;
}
