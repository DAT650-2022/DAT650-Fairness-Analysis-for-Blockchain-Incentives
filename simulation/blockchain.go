package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrTxNotFound    = errors.New("transaction not found")
	ErrNoValidTx     = errors.New("there is no valid transaction")
	ErrBlockNotFound = errors.New("block not found")
	ErrInvalidBlock  = errors.New("block is not valid")
)

// Blockchain keeps a sequence of Blocks
type Blockchain struct {
	blocks []*Block
}

// NewBlockchain creates a new blockchain with genesis Block
func NewBlockchain(address string) *Blockchain {
	// TODO(student)
	//tx := Transaction{}
	tx := NewCoinbaseTX(GenesisCoinbaseData, "")
	block := NewGenesisBlock(time.Now().Unix(), tx)
	var blocks []*Block
	blocks = append(blocks, block)
	return &Blockchain{blocks: blocks}
}

// addBlock saves the block into the blockchain
func (bc *Blockchain) addBlock(block *Block) error {
	// TODO(student) -- make sure you only add valid blocks!
	if !bc.ValidateBlock(block) {
		return ErrInvalidBlock
	}
	bc.blocks = append(bc.blocks, block)
	return nil
}

// GetGenesisBlock returns the Genesis Block
func (bc Blockchain) GetGenesisBlock() *Block {
	// TODO(student) -- YOU DON'T NEED TO CHANGE YOUR PREVIOUS METHOD
	return bc.blocks[0]
}

// CurrentBlock returns the last block
func (bc Blockchain) CurrentBlock() *Block {
	// TODO(student) -- YOU DON'T NEED TO CHANGE YOUR PREVIOUS METHOD
	return bc.blocks[len(bc.blocks)-1]
}

// GetBlock returns the block of a given hash
func (bc Blockchain) GetBlock(hash []byte) (*Block, error) {
	// TODO(student) -- YOU DON'T NEED TO CHANGE YOUR PREVIOUS METHOD
	for i, j := range bc.blocks {
		if bytes.Equal(j.Hash, hash) {
			return bc.blocks[i], nil
		}
	}
	return nil, ErrBlockNotFound
}

// ValidateBlock validates the block before adding it to the blockchain
func (bc *Blockchain) ValidateBlock(block *Block) bool {
	// TODO(student) -- a valid block cannot be nil and must contain txs.
	// Also, it should has the result of a valid PoW.
	if block != nil {
		if len(block.Transactions) > 0 {
			if block.Hash != nil {
				pow := NewProofOfWork(block)
				return pow.Validate()
			}
		}
	}
	return false
}

// MineBlock mines a new block with the provided transactions
func (bc *Blockchain) MineBlock(transactions []*Transaction) (*Block, error) {
	// TODO(student)
	// 1) Verify the existence of transactions inputs and discard invalid transactions that make reference to unknown inputs
	// 2) Add a block if there is a list of valid transactions
	if len(transactions) == 0 {
		return nil, ErrNoValidTx
	}
	for _, transaction := range transactions {
		if !bc.VerifyTransaction(transaction) {
			return nil, ErrNoValidTx
		}
	}
	block := NewBlock(time.Now().Unix(), transactions, bc.CurrentBlock().Hash)
	block.Mine()
	if bc.ValidateBlock(block) {
		bc.addBlock(block)
		return block, nil
	}
	return nil, ErrNoValidTx
}

// VerifyTransaction verifies if referred inputs exist
func (bc Blockchain) VerifyTransaction(tx *Transaction) bool {
	// TODO(student)
	// Check if all inputs of a given transaction refer to a existent transaction made previously
	// if not, you should return false!
	// TIP: remember that Coinbase transaction doesn't have input. Thus all coinbase tx are valid
	// if its coinbase then its valid
	if tx.IsCoinbase() {
		return true
	}
	// range through the inputs to make sure they're all valid
	for _, in := range tx.Vin {
		_, err := bc.FindTransaction(in.Txid)
		if err != nil {
			return false
		}
	}
	return true
}

// FindTransaction finds a transaction by its ID in the whole blockchain
func (bc Blockchain) FindTransaction(ID []byte) (*Transaction, error) {
	// TODO(student)
	// TIP: the chain is made of what?
	for _, i := range bc.blocks {
		for _, j := range i.Transactions {
			if bytes.Equal(j.ID, ID) {
				return j, nil
			}
		}
	}
	return nil, errors.New("Transaction not found in any block")
}

// FindUTXOSet finds and returns all unspent transaction outputs
func (bc Blockchain) FindUTXOSet() UTXOSet {
	// TODO(student)
	// 1) Search in the blockchain for unspent transactions outputs
	// 2) Ignore an already spent output
	// TIP: what determines that an output was spent?
	// create the utxo set
	utxo := make(UTXOSet)
	// range through all the blocks in the blockchain
	for _, block := range bc.blocks {
		// range through the transactions of the block
		for _, transaction := range block.Transactions {
			// create the utxo submap for the transaction
			utxo[Bytes2Hex(transaction.ID)] = make(map[int]TXOutput)
			// add all the outputs of the current transaction
			for i, out := range transaction.Vout {
				utxo[Bytes2Hex(transaction.ID)][i] = out
			}
			// check all the inputs of the current transaction to delete spent outputs
			// alternatively could've probably done this first to avoid adding things and then deleting them but w/e
			for _, in := range transaction.Vin {
				_, present := utxo[Bytes2Hex(in.Txid)][in.OutIdx]
				if present {
					delete(utxo[Bytes2Hex(in.Txid)], in.OutIdx)
					if len(utxo[Bytes2Hex(in.Txid)]) == 0 {
						delete(utxo, Bytes2Hex(in.Txid))
					}
				}
			}
		}
	}
	return utxo
}

func (bc Blockchain) String() string {
	var lines []string
	for _, block := range bc.blocks {
		lines = append(lines, fmt.Sprintf("%v", block))
	}
	return strings.Join(lines, "\n")
}
