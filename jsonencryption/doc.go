// Package jsonencryption provides a simple way to handle encryption and decryption
// of JSON data, particularly useful for database interactions. It introduces
// EncryptedID, a type designed to replace auto-incrementing primary keys in databases
// with an encrypted form in JSON, enhancing data privacy and security.
//
// Key Features:
// - EncryptedID: A struct type to represent encrypted database primary keys.
// - Encryption and Decryption: Functions to convert between encrypted and plain data.
// - Hashing and Masking: Tools to hash and mask uint values for additional security.
//
// Usage involves setting a 16-byte encryption key and utilizing the provided
// methods to encrypt JSON data structures, primarily focusing on the primary key fields.
// EncryptedID can be seamlessly integrated into database models and JSON serialization,
// ensuring that sensitive integer IDs are not exposed in API responses.
//
// The package also includes functions for hashing and masking uints, useful for
// creating non-reversible identifiers or temporary masked values for secure data handling.
package jsonencryption
