package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"

	"golang.org/x/crypto/ripemd160"
)

const (
	version            = byte(0x00)
	addressChecksumLen = 4
)

// newKeyPair creates a new cryptographic key pair
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	// TODO(student)
	// Create a new cryptographic key pair using the "elliptic" and "ecdsa" package.
	// Additionally, convert the PublicKey to bytes and return it (look at pubKeyToByte function).
	privkey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	pubkey := privkey.PublicKey
	pubkeybyt := pubKeyToByte(pubkey)
	return *privkey, pubkeybyt
}

// pubKeyToByte converts the ecdsa.PublicKey to a concatenation of its coordinates in bytes
func pubKeyToByte(pubkey ecdsa.PublicKey) []byte {
	// TODO(student)
	// step 1 of: https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses#How_to_create_Bitcoin_Address
	// 1: Take the corresponding public key generated with it (33 bytes, 1 byte 0x02 (y-coord is even), and 32 bytes corresponding to X coordinate)
	y := pubkey.Y.Bytes()
	x := pubkey.X.Bytes()
	pubByt := append(x, y...)
	return pubByt
}

// GetAddress returns address
// https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses#How_to_create_Bitcoin_Address
func GetAddress(pubKeyBytes []byte) []byte {
	// TODO(student)
	// Create a address following the logic described in the link above and
	// in the lab documentation.
	// 2 & 3
	hash := HashPubKey(pubKeyBytes)
	// 4: Add version byte in front of RIPEMD-160 hash (0x00 for Main Network)
	hashver := append([]byte{version}, hash...)
	// 5 & 6 & 7
	check := checksum(hashver)
	// 8: Add the 4 checksum bytes from stage 7 at the end of extended RIPEMD-160 hash from stage 4. This is the 25-byte binary Bitcoin Address.
	add := append(hashver, check...)
	// 9: Convert the result from a byte string into a base58 string using Base58Check encoding. This is the most commonly used Bitcoin Address format
	encoded := Base58Encode(add)
	return encoded
}

// GetStringAddress returns address as string
func GetStringAddress(pubKeyBytes []byte) string {
	// 2 & 3
	hash := HashPubKey(pubKeyBytes)
	// 4: Add version byte in front of RIPEMD-160 hash (0x00 for Main Network)
	hashver := append([]byte{version}, hash...)
	// 5 & 6 & 7
	check := checksum(hashver)
	// 8: Add the 4 checksum bytes from stage 7 at the end of extended RIPEMD-160 hash from stage 4. This is the 25-byte binary Bitcoin Address.
	add := append(hashver, check...)
	// 9: Convert the result from a byte string into a base58 string using Base58Check encoding. This is the most commonly used Bitcoin Address format
	encoded := Base58Encode(add)
	return string(encoded)
}

// HashPubKey hashes public key
func HashPubKey(pubKey []byte) []byte {
	// TODO(student)
	// compute the SHA256 + RIPEMD160 hash of the pubkey
	// step 2 and 3 of: https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses#How_to_create_Bitcoin_Address
	// use the go package ripemd160:
	// https://godoc.org/golang.org/x/crypto/ripemd160
	// 2 - Perform SHA-256 hashing on the public key
	hash256 := sha256.Sum256(pubKey)
	// 3 - Perform RIPEMD-160 hashing on the result of SHA-256
	hash160 := ripemd160.New()
	hash160.Write(hash256[:])
	return hash160.Sum(nil)
}

// GetPubKeyHashFromAddress returns the hash of the public key
// discarding the version and the checksum
func GetPubKeyHashFromAddress(address string) []byte {
	// TODO(student)
	// Decode the address using Base58Decode and extract the hash of the pubkey
	// Look in the picture of the documentation of the lab to understand
	// how it is stored: version + pubkeyhash + checksum
	addbytes := []byte(address)
	dec := Base58Decode(addbytes)
	// exclude the version and checksum
	return dec[1 : len(dec)-4]
}

// ValidateAddress check if an address is valid
func ValidateAddress(address string) bool {
	// TODO(student)
	// Validate a address by decoding it, extracting the
	// checksum, re-computing it using the "checksum" function
	// and comparing both.
	addbytes := []byte(address)
	dec := Base58Decode(addbytes)
	checkorig := dec[len(dec)-4:]
	check := checksum(dec[:len(dec)-4])
	return bytes.Equal(checkorig, check)
}

// Checksum generates a checksum for a public key
func checksum(payload []byte) []byte {
	// TODO(student)
	// Perform a double sha256 on the versioned payload
	// and return the first 4 bytes
	// Steps 5,6, and 7 of: https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses#How_to_create_Bitcoin_Address
	// 5 - Perform SHA-256 hash on the extended RIPEMD-160 result
	sha := sha256.Sum256(payload)
	// 6 - Perform SHA-256 hash on the result of the previous SHA-256 hash
	sha256again := sha256.Sum256(sha[:])
	// 7 - Take the first 4 bytes of the second SHA-256 hash. This is the address checksum
	return sha256again[:4]
}

func encodeKeyPair(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (string, string) {
	return encodePrivateKey(privateKey), encodePublicKey(publicKey)
}

func encodePrivateKey(privateKey *ecdsa.PrivateKey) string {
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	return string(pemEncoded)
}

func encodePublicKey(publicKey *ecdsa.PublicKey) string {
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return string(pemEncodedPub)
}

func decodeKeyPair(pemEncoded string, pemEncodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	return decodePrivateKey(pemEncoded), decodePublicKey(pemEncodedPub)
}

func decodePrivateKey(pemEncoded string) *ecdsa.PrivateKey {
	block, _ := pem.Decode([]byte(pemEncoded))
	privateKey, _ := x509.ParseECPrivateKey(block.Bytes)

	return privateKey
}

func decodePublicKey(pemEncodedPub string) *ecdsa.PublicKey {
	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	genericPubKey, _ := x509.ParsePKIXPublicKey(blockPub.Bytes)
	publicKey := genericPubKey.(*ecdsa.PublicKey) // cast to ecdsa

	return publicKey
}
