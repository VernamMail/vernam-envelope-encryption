// encrypt_decrypt is a minimal end-to-end example of the envelope-field
// encryption portion of the protocol.
//
// Run from the module root (go/):
//
//	go run ./examples/encrypt_decrypt
package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/vernammail/envelope-encryption"
)

func main() {
	// Step 1: generate a fresh per-email session key.
	key, err := envelope.NewSessionKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "session key generation failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Session key generated (%d bytes).\n", len(key))

	// Step 2: encrypt several envelope fields independently with the same key.
	// Each EncryptField call uses a fresh random IV internally.
	fields := map[string]string{
		"from":    "alice@example.com",
		"to":      "bob@example.org",
		"subject": "Quarterly report draft",
	}

	encrypted := make(map[string]envelope.EncryptedField, len(fields))
	for name, plaintext := range fields {
		ct, err := envelope.EncryptField(key, []byte(plaintext))
		if err != nil {
			fmt.Fprintf(os.Stderr, "EncryptField(%q) failed: %v\n", name, err)
			os.Exit(1)
		}
		encrypted[name] = ct
		fmt.Printf("Encrypted %-7s: %d bytes (IV+ciphertext+tag), hex prefix %s...\n",
			name, len(ct), hex.EncodeToString(ct[:8]))
	}

	// Step 3: decrypt each field back to plaintext.
	fmt.Println()
	for name, ct := range encrypted {
		pt, err := envelope.DecryptField(key, ct)
		if err != nil {
			fmt.Fprintf(os.Stderr, "DecryptField(%q) failed: %v\n", name, err)
			os.Exit(1)
		}
		fmt.Printf("Decrypted %-7s: %q\n", name, string(pt))
	}

	// Step 4: demonstrate that AEAD authentication rejects tampering.
	tampered := make(envelope.EncryptedField, len(encrypted["subject"]))
	copy(tampered, encrypted["subject"])
	tampered[len(tampered)-1] ^= 0x01

	fmt.Println()
	if _, err := envelope.DecryptField(key, tampered); err != nil {
		fmt.Printf("Tampered subject rejected as expected: %v\n", err)
	} else {
		fmt.Fprintln(os.Stderr, "ERROR: tampered ciphertext was accepted; this should never happen")
		os.Exit(1)
	}
}
