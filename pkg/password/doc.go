// Package password implements functions for securly storing user passwords.
//
// Security
//
// A plaintext password secured with this package is protected by multiple
// layers:
//
// 1. The plaintext password is transformed into a hash value using SHA3-512.
// This addresses two particular issues with bcrypt: Some implementations of
// bcrypt truncate the input to 72 bytes, which reduces the entropy of the
// passwords. Other implementations don’t truncate the input and are therefore
// vulnerable to DoS attacks because they allow the input of arbitrarily long
// passwords. By applying SHA3, a really long password can be quickly converted
// into a fixed length 512 bit value, solving both problems. The hash is than
// base64 encoded because bcrypt will stop at a null byte.
//
// 2. The base64 encoded SHA3-512 hash is hashed again using bcrypt with a
// specified cost (should be at least 10). Unlike cryptographic hash functions
// like SHA3, bcrypt is designed to be slow and hard to speed up via custom
// hardware and GPUs. A work factor of 10 translates into roughly 100ms for all
// these steps.
//
// 3. The resulting bcrypt hash is encrypted with AES (AES-128, AES-192, AES-256
// are selectable by choosing the appropriate secret key length) using a secret
// key common to all hashes and referred to as a pepper. The pepper is a defense
// in depth measure. Its value should be stored separately in a manner that
// makes it difficult to discover by an attacker (i.e. not in a database table).
// As a result, if only the password storage is compromised, the password hashes
// are encrypted and of no use to an attacker.
//
// Why not use scrypt or argon2 over bcrypt
//
// Most security experts agree that scrypt and bcrypt provide similar
// protections.
//
// While argon2 is a fantastic password hashing function, bcrypt has been around
// since 1999 without any significant vulnerabilities found.
//
// Why is the pepper used for encryption instead of hashing
//
// Recall that the pepper is a defense in depth measure and should be stored
// separately. But storing it separately also includes the possibility of the
// pepper (and not the password hashes) being compromised. If the pepper is used
// for hashing, it cannot be easily rotated. Instead, using it for encryption
// gives similar security but with the added ability to rotate.
package password
