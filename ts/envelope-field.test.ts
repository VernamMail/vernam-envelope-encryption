// Tests for the TypeScript envelope-field reference implementation.
//
// Run with:  npm test   (uses Node's built-in test runner via tsx)
//
// The cross-language test (TestKnownVectorMatchesGo) is the important one:
// it asserts that this TypeScript implementation produces byte-identical
// output to the Go reference for the shared test vector envelope-field-001.
// This is the concrete proof that the specification is language-agnostic.

import { test } from "node:test";
import assert from "node:assert/strict";
import {
  newSessionKey,
  encryptField,
  encryptFieldWithIV,
  decryptField,
  toHex,
  fromHex,
  SESSION_KEY_SIZE,
  GCM_NONCE_SIZE,
  GCM_TAG_SIZE,
  AuthenticationError,
  MalformedWireError,
} from "./envelope-field.ts";

const enc = new TextEncoder();
const dec = new TextDecoder();

test("newSessionKey returns 32 distinct random bytes", () => {
  const a = newSessionKey();
  const b = newSessionKey();
  assert.equal(a.length, SESSION_KEY_SIZE);
  assert.equal(a.length, 32);
  assert.notDeepEqual([...a], [...b], "two session keys should differ");
});

test("encryptField round-trips", async () => {
  const key = newSessionKey();
  const cases = [
    "alice@example.com",
    "bob+tag@subdomain.example.org",
    "",
    'From: "Quoted Name" <name@example.com>, second@example.com',
  ];
  for (const pt of cases) {
    const ct = await encryptField(key, enc.encode(pt));
    const got = dec.decode(await decryptField(key, ct));
    assert.equal(got, pt);
    assert.equal(ct.length, enc.encode(pt).length + GCM_NONCE_SIZE + GCM_TAG_SIZE);
  }
});

// The cross-language interop proof. These values come from
// test-vectors/basic.json (envelope-field-001) and are produced by the Go
// reference in go/envelope_test.go::TestEncryptFieldKnownVector.
test("TestKnownVectorMatchesGo", async () => {
  const keyHex =
    "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f";
  const ivHex = "000102030405060708090a0b";
  const plaintext = "alice@example.com";
  const expectedWireHex =
    "000102030405060708090a0b266ebf78a0a5a763ec2ce7e7d4c71b02eeb13438bcf297aba04675c1f322b9610e";

  const key = fromHex(keyHex);
  const iv = fromHex(ivHex);
  const got = await encryptFieldWithIV(key, iv, enc.encode(plaintext));

  assert.equal(
    toHex(got),
    expectedWireHex,
    "TS wire bytes must match the Go reference exactly",
  );

  const roundTripped = dec.decode(await decryptField(key, got));
  assert.equal(roundTripped, plaintext);
});

test("decryptField rejects tampering", async () => {
  const key = newSessionKey();
  const ct = await encryptField(key, enc.encode("alice@example.com"));
  const tampered = new Uint8Array(ct);
  tampered[tampered.length - 1] ^= 0x01;
  await assert.rejects(() => decryptField(key, tampered), AuthenticationError);
});

test("decryptField rejects too-short input", async () => {
  const key = newSessionKey();
  for (const len of [0, 11, 27]) {
    await assert.rejects(
      () => decryptField(key, new Uint8Array(len)),
      MalformedWireError,
    );
  }
});
