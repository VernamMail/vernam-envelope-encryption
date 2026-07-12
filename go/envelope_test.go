package envelope

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// TestSessionKeyGenerationDistinct verifies that successive session keys
// are not identical (a basic CSPRNG sanity check).
func TestSessionKeyGenerationDistinct(t *testing.T) {
	k1, err := NewSessionKey()
	if err != nil {
		t.Fatalf("NewSessionKey returned error: %v", err)
	}
	k2, err := NewSessionKey()
	if err != nil {
		t.Fatalf("NewSessionKey returned error: %v", err)
	}
	if bytes.Equal(k1[:], k2[:]) {
		t.Fatal("two successive session keys were identical; CSPRNG appears broken")
	}
}

// TestSessionKeySize verifies the session key is exactly 32 bytes.
func TestSessionKeySize(t *testing.T) {
	k, err := NewSessionKey()
	if err != nil {
		t.Fatalf("NewSessionKey returned error: %v", err)
	}
	if len(k) != SessionKeySize {
		t.Fatalf("session key size = %d, want %d", len(k), SessionKeySize)
	}
	if SessionKeySize != 32 {
		t.Fatalf("SessionKeySize = %d, want 32 per SPEC.md §4.2", SessionKeySize)
	}
}

// TestWrappedKeyWireSize verifies the wire-format constant matches the
// specification (SPEC.md §5.1: 1 + 1568 + 32 + 12 + 32 + 16 = 1661).
func TestWrappedKeyWireSize(t *testing.T) {
	expected := 1 + KyberCiphertextSize + X25519KeySize + GCMNonceSize + SessionKeySize + GCMTagSize
	if WrappedKeyWireSize != expected {
		t.Fatalf("WrappedKeyWireSize = %d, want %d (per SPEC.md §5.1)", WrappedKeyWireSize, expected)
	}
	if WrappedKeyWireSize != 1661 {
		t.Fatalf("WrappedKeyWireSize = %d, want 1661 per SPEC.md §5.1", WrappedKeyWireSize)
	}
}

// TestHKDFInfoString verifies the HKDF info string is exactly the value
// required by the spec for v0.1 compliance (SPEC.md §4.3 step 4).
func TestHKDFInfoString(t *testing.T) {
	if HKDFInfo != "EnvelopeMetadataEncryption-v1" {
		t.Fatalf("HKDFInfo = %q, want %q (per SPEC.md §4.3)", HKDFInfo, "EnvelopeMetadataEncryption-v1")
	}
}

// TestVersionConstants verifies the wire version constants.
func TestVersionConstants(t *testing.T) {
	if WireVersionV1 != 0x01 {
		t.Fatalf("WireVersionV1 = 0x%02x, want 0x01", WireVersionV1)
	}
	if WireVersionV2 != 0x02 {
		t.Fatalf("WireVersionV2 = 0x%02x, want 0x02", WireVersionV2)
	}
}

// TestEncryptFieldRoundTrip verifies that production encryption (random IV)
// followed by decryption recovers the original plaintext.
func TestEncryptFieldRoundTrip(t *testing.T) {
	key, err := NewSessionKey()
	if err != nil {
		t.Fatalf("NewSessionKey: %v", err)
	}
	cases := []string{
		"alice@example.com",
		"bob+tag@subdomain.example.org",
		"",
		"From: \"Quoted Name\" <name@example.com>, second@example.com",
	}
	for _, plaintext := range cases {
		ct, err := EncryptField(key, []byte(plaintext))
		if err != nil {
			t.Fatalf("EncryptField(%q): %v", plaintext, err)
		}
		got, err := DecryptField(key, ct)
		if err != nil {
			t.Fatalf("DecryptField(%q): %v", plaintext, err)
		}
		if string(got) != plaintext {
			t.Fatalf("round-trip mismatch: got %q, want %q", string(got), plaintext)
		}
		expectedLen := len(plaintext) + GCMNonceSize + GCMTagSize
		if len(ct) != expectedLen {
			t.Fatalf("ciphertext length = %d, want %d (per SPEC.md §5.2)", len(ct), expectedLen)
		}
	}
}

// TestEncryptFieldKnownVector verifies that the deterministic core of
// envelope-field encryption produces exactly the bytes recorded in
// test-vectors/basic.json (vector envelope-field-001). This catches any
// regression that would silently change the wire format and break
// interoperability with previously-encrypted data.
func TestEncryptFieldKnownVector(t *testing.T) {
	keyHex := "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"
	ivHex := "000102030405060708090a0b"
	plaintext := "alice@example.com"
	expectedWireHex := "000102030405060708090a0b266ebf78a0a5a763ec2ce7e7d4c71b02eeb13438bcf297aba04675c1f322b9610e"

	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil || len(keyBytes) != SessionKeySize {
		t.Fatalf("bad test key")
	}
	var key SessionKey
	copy(key[:], keyBytes)
	iv, err := hex.DecodeString(ivHex)
	if err != nil || len(iv) != GCMNonceSize {
		t.Fatalf("bad test IV")
	}
	expectedWire, err := hex.DecodeString(expectedWireHex)
	if err != nil {
		t.Fatalf("bad expected wire hex")
	}

	got, err := encryptFieldWithIV(key, iv, []byte(plaintext))
	if err != nil {
		t.Fatalf("encryptFieldWithIV: %v", err)
	}
	if !bytes.Equal(got, expectedWire) {
		t.Fatalf("wire mismatch:\n got  %s\n want %s", hex.EncodeToString(got), expectedWireHex)
	}

	roundTripped, err := DecryptField(key, got)
	if err != nil {
		t.Fatalf("DecryptField: %v", err)
	}
	if string(roundTripped) != plaintext {
		t.Fatalf("decrypt mismatch: got %q, want %q", string(roundTripped), plaintext)
	}
}

// TestDecryptFieldRejectsTamper verifies that AES-GCM authentication
// rejects any modification to the ciphertext or tag.
func TestDecryptFieldRejectsTamper(t *testing.T) {
	key, err := NewSessionKey()
	if err != nil {
		t.Fatalf("NewSessionKey: %v", err)
	}
	ct, err := EncryptField(key, []byte("alice@example.com"))
	if err != nil {
		t.Fatalf("EncryptField: %v", err)
	}
	tampered := make(EncryptedField, len(ct))
	copy(tampered, ct)
	tampered[len(tampered)-1] ^= 0x01

	if _, err := DecryptField(key, tampered); err != ErrAuthenticationFail {
		t.Fatalf("DecryptField on tampered input: got err=%v, want ErrAuthenticationFail", err)
	}
}

// TestDecryptFieldRejectsTooShort verifies that malformed inputs are rejected
// without attempting cryptographic operations.
func TestDecryptFieldRejectsTooShort(t *testing.T) {
	key, err := NewSessionKey()
	if err != nil {
		t.Fatalf("NewSessionKey: %v", err)
	}
	for _, length := range []int{0, 11, 27} { // less than IV+tag = 28
		_, err := DecryptField(key, make(EncryptedField, length))
		if err != ErrMalformedWire {
			t.Fatalf("DecryptField(len=%d): got err=%v, want ErrMalformedWire", length, err)
		}
	}
}

// TestWrapNotYetImplemented documents that hybrid wrap is scheduled for
// milestone 3 of the NLnet-funded work.
func TestWrapNotYetImplemented(t *testing.T) {
	k, err := NewSessionKey()
	if err != nil {
		t.Fatalf("NewSessionKey: %v", err)
	}
	_, err = Wrap(k, PublicIdentity{})
	if err == nil {
		t.Skip("Wrap is now implemented; remove this guard test and add real round-trip tests")
	}
}
