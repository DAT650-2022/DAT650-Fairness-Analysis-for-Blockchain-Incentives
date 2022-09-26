package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
)

var (
	ErrNoFunds         = errors.New("not enough funds")
	ErrTxInputNotFound = errors.New("transaction input not found")
)

// Transaction represents a Bitcoin transaction
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to, data string) (*Transaction, error) {
	// TODO(student) -- update your function to generate random data
	// if the data is not provided
	if data == "" {
		data = fmt.Sprintf("Reward to %s", to)
	}

	// vin := &TXInput{Txid: nil, OutIdx: -1, ScriptSig: data}
	vin := &TXInput{Txid: nil, OutIdx: -1, PubKey: []byte(data), Signature: nil}
	// vout := &TXOutput{Value: BlockReward, ScriptPubKey: to}
	vout := NewTXOutput(BlockReward, to)

	var vinSlice = []TXInput{*vin}
	var voutSlice = []TXOutput{*vout}

	tx := &Transaction{
		ID:   nil,
		Vin:  vinSlice,
		Vout: voutSlice,
	}
	tx.Hash()
	return tx, nil
}

// NewUTXOTransaction creates a new UTXO transaction
// NOTE: The returned tx is NOT signed!
func NewUTXOTransaction(pubKey []byte, to string, amount int, utxos UTXOSet) (*Transaction, error) {
	// TODO(student)
	// Modify your function to use the address instead of just strings
	pubKeyHash := HashPubKey(pubKey)
	money, spendutxo := utxos.FindSpendableOutputs(pubKeyHash, amount)

	if money < amount {
		return nil, ErrNoFunds
	}
	var ins []TXInput
	// range through the spendable utxos
	for txid, outIdxs := range spendutxo {
		// range through the out indexes
		for _, outIdx := range outIdxs {
			// add input for each used out
			input := TXInput{
				Txid:      Hex2Bytes(txid),
				OutIdx:    outIdx,
				Signature: nil,
				PubKey:    pubKey,
			}
			ins = append(ins, input)
		}
	}
	// create the transaction and add the inputs
	tx := Transaction{Vin: ins}
	// create he outputs and add the recipient of the transaction
	out := NewTXOutput(amount, to)
	outs := []TXOutput{*out}
	// if theres change left over create an extra output
	if money > amount {
		diff := money - amount
		// need to get the address of the sender
		from := GetAddress(pubKey)
		// use the new create output function instead of making it raw
		change := NewTXOutput(diff, string(from))
		outs = append(outs, *change)
	}
	// add the outputs to the transaction
	tx.Vout = outs
	// hash the id of the transaction
	tx.ID = tx.Hash()
	return &tx, nil
}

// IsCoinbase checks whether the transaction is coinbase
func (tx Transaction) IsCoinbase() bool {
	// TODO(student)
	return tx.Vin[0].OutIdx == -1 && tx.Vin[0].Txid == nil
}

// Equals checks if the given transaction ID matches the ID of tx
func (tx Transaction) Equals(ID []byte) bool {
	return bytes.Equal(tx.ID, ID)
}

// Serialize returns a serialized Transaction
func (tx Transaction) Serialize() []byte {
	// TODO(student) -- YOU DON'T NEED TO CHANGE YOUR PREVIOUS METHOD
	var b bytes.Buffer
	e := gob.NewEncoder(&b)
	err := e.Encode(tx)
	if err != nil {
		fmt.Println("Failed gob encode transaction serialize")
		panic(err)
	}
	return b.Bytes()
}

// Hash returns the hash of the Transaction
func (tx *Transaction) Hash() []byte {
	// TODO(student) -- YOU DON'T NEED TO CHANGE YOUR PREVIOUS METHOD
	txCopy := tx
	txCopy.ID = nil

	var b bytes.Buffer
	e := gob.NewEncoder(&b)
	err := e.Encode(txCopy)
	if err != nil {
		fmt.Println("Failed gob encode transaction hash")
		panic(err)
	}
	hashed := sha256.Sum256(b.Bytes())
	tx.ID = hashed[:]
	return hashed[:]
}

// TrimmedCopy creates a trimmed copy of Transaction to be used in signing
func (tx Transaction) TrimmedCopy() Transaction {
	// TODO(student)
	// You need to create a copy of the transaction to be signed.
	// The fields Signature and PubKey of the input need to be nil
	// since they are not included in signature.
	var ins []TXInput
	var outs []TXOutput
	// range through the inputs
	for _, in := range tx.Vin {
		// add all the inputs
		ins = append(ins, TXInput{Txid: in.Txid, OutIdx: in.OutIdx})
	}
	// same for outputs
	for _, out := range tx.Vout {
		outs = append(outs, TXOutput{out.Value, out.PubKeyHash})
	}
	// return the transaction with all the inputs and outputs
	return Transaction{tx.ID, ins, outs}
}

// Sign signs each input of a Transaction
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]*Transaction) error {
	// TODO(student)
	// 1) coinbase transactions are not signed.
	// 2) Throw a Panic in case of any prevTXs (used inputs) didn't exists
	// Take a look on the tests to see the expected error message
	// 3) Create a copy of the transaction to be signed
	// 4) Sign all the previous TXInputs of the transaction tx using the
	// copy as the payload (serialized) to be signed in the ecdsa.Sig
	// (https://golang.org/pkg/crypto/ecdsa/#Sign)
	// Make sure that each input of the copy to be signed
	// have the correct PubKeyHash of each output in the prevTXs
	// Store the signature as a concatenation of R and S fields
	// 1 coinbase is not signed
	if tx.IsCoinbase() {
		return nil
	}
	// 2 range through the inputs and make sure they have a txid
	for _, in := range tx.Vin {
		if prevTXs[Bytes2Hex(in.Txid)] == nil {
			return ErrTxInputNotFound
		}
	}
	// 3 create a copy of the tx
	txCopy := tx.TrimmedCopy()
	// 4 sign all previous tx
	for inID, vin := range txCopy.Vin {
		prevTx := prevTXs[Bytes2Hex(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.OutIdx].PubKeyHash

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.Serialize())
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vin[inID].Signature = signature
		txCopy.Vin[inID].PubKey = nil
	}
	return nil
}

// Verify verifies signatures of Transaction inputs
func (tx Transaction) Verify(prevTXs map[string]*Transaction) bool {
	// TODO(student)
	// 1) coinbase transactions are not signed.
	// 2) Throw a Panic in case of any prevTXs (used inputs) didn't exists
	// Take a look on the tests to see the expected error message
	// 3) Create the same copy of the transaction that was signed
	// and get the curve used for sign: P256
	// 4) Doing the opposite operation of the signing, perform the
	// verification of the signature, by recovering the R and S byte fields
	// of the Signature and the X and Y fields of the PubKey from
	// the inputs of tx. Verify the signature of each input using the
	// ecdsa.Verify function (https://golang.org/pkg/crypto/ecdsa/#Verify)
	// Note that to use this function you need to reconstruct the
	// ecdsa.PublicKey. Also notice that the ecdsa.Verify function receive
	// a byte array, you the transaction copy need to be serialized.
	// return true if all inputs have valid signature,
	// and false if any of them have an invalid signature.
	// 1 coinbase is not signed
	if tx.IsCoinbase() {
		return true
	}
	// 2 panic if prevtx doesnt exist | dont panic!
	for _, in := range tx.Vin {
		if prevTXs[Bytes2Hex(in.Txid)] == nil {
			return false
		}
	}
	// 3 create a copy and get a curve for the sign
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()
	// 4
	for id, in := range tx.Vin {
		prevTx := prevTXs[Bytes2Hex(in.Txid)]
		txCopy.Vin[id].PubKey = prevTx.Vout[in.OutIdx].PubKeyHash
		txCopy.ID = txCopy.Hash()

		// reconstruct R & S
		r := big.Int{}
		s := big.Int{}
		sigLen := len(in.Signature)
		r.SetBytes(in.Signature[:(sigLen / 2)])
		s.SetBytes(in.Signature[(sigLen / 2):])

		// reconstruct X & Y
		x := big.Int{}
		y := big.Int{}
		keyLen := len(in.PubKey)
		x.SetBytes(in.PubKey[:(keyLen / 2)])
		y.SetBytes(in.PubKey[(keyLen / 2):])
		rawPubKey := ecdsa.PublicKey{curve, &x, &y}

		return ecdsa.Verify(&rawPubKey, txCopy.Serialize(), &r, &s)
	}

	return false
}

// String returns a human-readable representation of a transaction
func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x :", tx.ID))

	for i, input := range tx.Vin {
		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Txid))
		lines = append(lines, fmt.Sprintf("       OutIdx:    %d", input.OutIdx))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey: %x", input.PubKey))
	}

	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       PubKeyHash: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}
