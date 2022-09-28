package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
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
func NewCoinbaseTX(to, data string) *Transaction {
	// TODO(student)
	// Create a new coinbase using the given data field
	// or the default "fmt.Sprintf("Reward to %s", to)"
	// if data is empty.
	if data == "" {
		data = fmt.Sprintf("Reward to %s", to)
	}

	vin := &TXInput{Txid: nil, OutIdx: -1, ScriptSig: data}
	vout := &TXOutput{Value: BlockReward, ScriptPubKey: to}

	var vinSlice = []TXInput{*vin}
	var voutSlice = []TXOutput{*vout}

	tx := &Transaction{
		ID:   nil,
		Vin:  vinSlice,
		Vout: voutSlice,
	}
	tx.Hash()
	return tx
}

// NewUTXOTransaction creates a new UTXO transaction
func NewUTXOTransaction(from, to string, amount int, utxos UTXOSet) (*Transaction, error) {
	// TODO(student)
	// 1) Find valid spendable outputs and the current balance of the sender
	// 2) The sender has sufficient funds? If not return the error:
	// "Not enough funds"
	// 3) Build a list of inputs based on the current valid outputs
	// 4) Build a list of new outputs, creating a "change" output if necessary
	// 5) Create a new transaction with the input and output list.
	money, spendutxo := utxos.FindSpendableOutputs(from, amount)
	if money < amount {
		return nil, ErrNoFunds
	}
	var ins []TXInput
	// range through the spendable utxos
	for txid, outIdxs := range spendutxo {
		// range through the out indexes
		for _, outIdx := range outIdxs {
			// add input for each used out
			input := TXInput{}
			input.Txid = Hex2Bytes(txid)
			input.ScriptSig = from
			input.OutIdx = outIdx
			ins = append(ins, input)
		}
	}
	// create the transaction and add the inputs
	tx := Transaction{}
	tx.Vin = ins
	// create he outputs and add the recipient of the transaction
	out := TXOutput{Value: amount, ScriptPubKey: to}
	outs := []TXOutput{out}
	// if theres change left over create an extra output
	if money > amount {
		diff := money - amount
		change := TXOutput{Value: diff, ScriptPubKey: from}
		outs = append(outs, change)
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
	// TIP: What differentiate a coinbase transaction from a normal transaction?
	// Remember that OutIdx represents the position of an output referred by the input
	return tx.Vin[0].OutIdx == -1 && tx.Vin[0].Txid == nil
}

// Equals checks if the given transaction ID matches the ID of tx
func (tx Transaction) Equals(ID []byte) bool {
	// TODO(student)
	return bytes.Equal(tx.ID, ID)
}

// Serialize returns a serialized Transaction
func (tx Transaction) Serialize() []byte {
	// TODO(student)
	// This function should encode all fields of the Transaction struct, using the gob encoder
	// Note: This includes the tx.ID!
	// TIP: https://golang.org/pkg/encoding/gob/
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
	// TODO(student)
	// This function should hash the serialized representation of a transaction but it MUST
	// ignore the ID (set it to nil), since the ID is the hash of the tx itself (if exists).
	// You may need to make a copy of the object, otherwise it will change the original pointer.
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

// String returns a human-readable representation of a transaction
func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x :", tx.ID))

	for i, input := range tx.Vin {
		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Txid))
		lines = append(lines, fmt.Sprintf("       OutIdx:    %d", input.OutIdx))
		lines = append(lines, fmt.Sprintf("       ScriptSig: %s", input.ScriptSig))
	}

	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       ScriptPubKey: %s", output.ScriptPubKey))
	}

	return strings.Join(lines, "\n")
}
