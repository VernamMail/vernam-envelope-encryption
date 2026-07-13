// Package envelope implements the envelope metadata encryption protocol
// described in SPEC.md.
//
// Currently provides:
//   - Per-email session key generation (NewSessionKey)
//   - AES-256-GCM envelope-field encryption / decryption (EncryptField, DecryptField)
//   - Wire-format constants validated against the specification
//
// Planned as roadmap milestones (see ../ROADMAP.md):
//   - Hybrid ML-KEM-1024 + X25519 wrap/unwrap (Wrap, Unwrap)
//   - Comprehensive test vectors for the wrapped-key wire format
//   - Cross-language interop test suite
//   - Third-party cryptographic review
//
// See ../SPEC.md for the protocol specification.
package envelope

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
)

// Protocol constants. See SPEC.md §4 and §5.
const (
	WireVersionV1 byte = 0x01
	WireVersionV2 byte = 0x02

	SessionKeySize      = 32   // bytes
	GCMNonceSize        = 12   // bytes
	GCMTagSize          = 16   // bytes
	X25519KeySize       = 32   // bytes
	KyberCiphertextSize = 1568 // ML-KEM-1024 ciphertext
	WrappedKeyWireSize  = 1661 // V + kyber_ct + eph_pk + IV + wrapped + tag

	HKDFInfo = "EnvelopeMetadataEncryption-v1"
)

// Errors.
var (
	ErrUnsupportedVersion = errors.New("envelope: unsupported wire version")
	ErrMalformedWire      = errors.New("envelope: malformed wire format")
	ErrAuthenticationFail = errors.New("envelope: authentication failed")
	ErrInvalidNonceSize   = errors.New("envelope: invalid nonce size")
)

// SessionKey is a 32-byte symmetric key generated per email.
type SessionKey [SessionKeySize]byte

// PublicIdentity holds a recipient's public-key identity for wrapping.
type PublicIdentity struct {
	KyberPublicKey  []byte // ML-KEM-1024 public key
	X25519PublicKey [X25519KeySize]byte
}

// PrivateIdentity holds a recipient's private keys for unwrapping.
type PrivateIdentity struct {
	KyberPrivateKey  []byte // ML-KEM-1024 private key
	X25519PrivateKey [X25519KeySize]byte
}

// EncryptedField is the wire format for a single encrypted envelope field.
// Layout: IV (12 bytes) || Ciphertext (n bytes) || Tag (16 bytes).
type EncryptedField []byte

// WrappedKey is the 1661-byte wire format for a session key wrapped
// to a recipient via the hybrid KEM. See SPEC.md §5.1.
type WrappedKey [WrappedKeyWireSize]byte

// NewSessionKey generates a fresh 32-byte session key from the system CSPRNG.
func NewSessionKey() (SessionKey, error) {
	var k SessionKey
	if _, err := rand.Read(k[:]); err != nil {
		return SessionKey{}, err
	}
	return k, nil
}

// EncryptField encrypts a single envelope field under the given session key
// using AES-256-GCM with a fresh random IV. The output layout is defined
// in SPEC.md §5.2: IV (12 bytes) || ciphertext || tag (16 bytes).
func EncryptField(key SessionKey, plaintext []byte) (EncryptedField, error) {
	var iv [GCMNonceSize]byte
	if _, err := rand.Read(iv[:]); err != nil {
		return nil, err
	}
	return encryptFieldWithIV(key, iv[:], plaintext)
}

// encryptFieldWithIV is the deterministic core used by EncryptField and by
// test-vector generators. It is unexported because production callers MUST
// use a fresh random IV per encryption (a reused IV under the same key
// breaks AES-GCM's confidentiality and integrity).
func encryptFieldWithIV(key SessionKey, iv []byte, plaintext []byte) (EncryptedField, error) {
	if len(iv) != GCMNonceSize {
		return nil, ErrInvalidNonceSize
	}
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	out := make([]byte, GCMNonceSize, GCMNonceSize+len(plaintext)+GCMTagSize)
	copy(out, iv)
	out = aead.Seal(out, iv, plaintext, nil)
	return EncryptedField(out), nil
}

// DecryptField decrypts a previously encrypted envelope field.
// Returns ErrMalformedWire if the input is too short to contain IV + tag.
// Returns ErrAuthenticationFail if the GCM tag does not verify.
func DecryptField(key SessionKey, ciphertext EncryptedField) ([]byte, error) {
	if len(ciphertext) < GCMNonceSize+GCMTagSize {
		return nil, ErrMalformedWire
	}
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	iv := ciphertext[:GCMNonceSize]
	ct := ciphertext[GCMNonceSize:]
	plaintext, err := aead.Open(nil, iv, ct, nil)
	if err != nil {
		return nil, ErrAuthenticationFail
	}
	return plaintext, nil
}

// Wrap encapsulates a session key to the given recipient using the hybrid
// ML-KEM-1024 + X25519 construction with ciphertext-binding AAD (v2).
// See SPEC.md §4.3 and §5.1.
func Wrap(sk SessionKey, recipient PublicIdentity) (WrappedKey, error) {
	// TODO(nlnet-milestone-3): Implement hybrid encapsulation per SPEC.md §4.3:
	//   1. ML-KEM-1024 encapsulate to recipient.KyberPublicKey
	//   2. X25519 ephemeral keygen + DH with recipient.X25519PublicKey
	//   3. HKDF-SHA256 (extract-then-expand) over the combined shared secrets with info=HKDFInfo
	//   4. AES-256-GCM wrap session key with AAD = kyber_ct || eph_pk
	//   5. Serialize per §5.1 wire format
	return WrappedKey{}, errors.New("envelope: Wrap not yet implemented")
}

// Unwrap decapsulates a wrapped session key using the recipient's private
// identity. Returns ErrAuthenticationFail if the AAD-bound GCM tag does
// not verify, or if any constituent step fails.
func Unwrap(w WrappedKey, recipient PrivateIdentity) (SessionKey, error) {
	// TODO(nlnet-milestone-3): Implement hybrid decapsulation per SPEC.md §4.4.
	return SessionKey{}, errors.New("envelope: Unwrap not yet implemented")
}

// version returns the wire version byte of a wrapped key.
func (w WrappedKey) version() byte {
	return w[0]
}
