package main

import (
	"bytes"
	"fmt"
)

// TXOutput represents a transaction output
type TXOutput struct {
	Value      int    // The transaction value
	PubKeyHash []byte // The conditions to claim this output. For this demo we will use the hash of the public key (used to "lock" the output)
}

// Lock locks the transaction to a specific address
// Only this address owns this transaction
func (out *TXOutput) Lock(address string) {
	// TODO(student)
	// "Lock" the TXOutput to a specific PubKeyHash
	// based on the given address
	addByt := []byte(address)
	pubKeyHash := Base58Decode(addByt)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

// IsLockedWithKey checks if the output can be used by the owner of the pubkey
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	// TODO(student)
	return bytes.Equal(out.PubKeyHash, pubKeyHash)
}

// NewTXOutput create a new TXOutput
func NewTXOutput(value int, address string) *TXOutput {
	// TODO(student)
	// Create a new locked TXOutput
	out := TXOutput{
		Value: value,
	}
	out.Lock(address)
	return &out
}

func (out TXOutput) String() string {
	return fmt.Sprintf("{%d, %x}", out.Value, out.PubKeyHash)
}
