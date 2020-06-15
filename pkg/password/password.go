package password

import (
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"

	"github.com/my-cargonaut/cargonaut/pkg/crypto/aes"
)

// DefaultCost is a sensible recommendation for the cost which should be used.
const DefaultCost = bcrypt.DefaultCost + 2

// ErrPasswordMismatch indicates that the hash is not the hash of the provided
// plaintext password.
var ErrPasswordMismatch = errors.New("password: passwords do not match")

// Generate creates a secure password hash from the provided plaintext password.
// A cost between 4 (bcrypt.MinCost) and 31 (bcrypt.MaxCost) must be specified,
// the default cost is specified by the DefaultCost constant (12). The password
// hash emitted by bcrypt is additionally AES encrypted with the pepper. The
// pepper must be 16, 24 or 32 byte broad to select either AES-128, AES-192 or
// AES-256.
func Generate(pepper []byte, password string, cost int) (string, error) {
	// Step one: Hash the plaintext password using SHA3-512. This gives a fixed
	// length 64 byte value, even for arbitrarily long passwords. The resulting
	// hash is base64 encoded using standard encoding because bcrypt will stop
	// at a null byte.
	hashedPassword := sha3.Sum512([]byte(password))
	encodedPasswordHash := base64.StdEncoding.EncodeToString(hashedPassword[:])

	// Step two: Hash the password using bcrypt. The cost needs to be adjusted
	// as computing power rises.
	bcryptedPassword, err := bcrypt.GenerateFromPassword([]byte(encodedPasswordHash), cost)
	if err != nil {
		return "", err
	}

	// Step three: Encrypt the bcrypt hashed password with AES using the pepper
	// as encryption key.
	return aes.Encrypt(pepper, bcryptedPassword)
}

// Compare a plaintext password against a password hash created using bcrypt.
// Before this comparison can happen, the hash must be decrypted, since it is
// AES encrypted.
func Compare(pepper []byte, password, encryptedPasswordHash string) error {
	// Hash the plaintext password using SHA3-512 and base64 encode it using
	// standard encoding.
	hashedPassword := sha3.Sum512([]byte(password))
	encodedPasswordHash := base64.StdEncoding.EncodeToString(hashedPassword[:])

	// Decrypt the AES encrypted password hash to retrieve the hash which was
	// once created by bcrypt from the SHA3-512 checksum of the plaintext
	// password.
	bcryptedPassword, err := aes.Decrypt(pepper, encryptedPasswordHash)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrPasswordMismatch
	} else if err != nil {
		return err
	}

	// Compare the bcrypt password hash with the plaintext password.
	if err = bcrypt.CompareHashAndPassword([]byte(bcryptedPassword), []byte(encodedPasswordHash)); err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrPasswordMismatch
	}
	return err
}
