package main

// UTXOSet represents a set of UTXO as an in-memory cache
// The key of the most external map is the transaction ID
// (encoded as string) that contains these outputs
// {map of transaction ID -> {map of TXOutput Index -> TXOutput}}
type UTXOSet map[string]map[int]TXOutput

// FindSpendableOutputs finds and returns unspent outputs in the UTXO Set
// to reference in inputs and the current accumulated balance
func (u UTXOSet) FindSpendableOutputs(unlockingData string, amount int) (int, map[string][]int) {
	// TODO(student)
	unspent := make(map[string][]int)
	var accumulated int
	// iterates over the utxoset map
	for txid, txouts := range u {
		// iterates over the txouts of the current utxo
		for outIdx, txout := range txouts {
			if txout.CanBeUnlockedWith(unlockingData) {
				// if it can be unlocked then add it to the money and add the transaction id to the map
				accumulated += txout.Value
				unspent[txid] = append(unspent[txid], outIdx)
			}
			if accumulated > amount {
				return accumulated, unspent
			}
		}
	}
	return 0, make(map[string][]int)
}

// FindUTXO finds all UTXO in the UTXO Set for a given unlockingData key (e.g., address)
// This function ignores the index of each output and returns
// a list of all outputs in the UTXO Set that can be unlocked by the user
func (u UTXOSet) FindUTXO(unlockingData string) []TXOutput {
	var UTXO []TXOutput
	// TODO(student)
	// iterate over the map
	for _, j := range u {
		// go through txouts of current map entry
		for _, m := range j {
			// add it if it can be unlocked
			if m.CanBeUnlockedWith(unlockingData) {
				UTXO = append(UTXO, m)
			}
		}
	}
	return UTXO
}

// CountUTXOs returns the number of transactions outputs in the UTXO set
func (u UTXOSet) CountUTXOs() int {
	// TODO(student)
	var count int
	for _, j := range u {
		// add all txout of current map entry
		count += len(j)
	}
	return count
}

// Update updates the UTXO Set with the new set of transactions
func (u UTXOSet) Update(transactions []*Transaction) {
	// TODO(student)
	// Iterate over the transactions  and update
	// the current UTXOSet with the new
	// transactions.
	//
	// TIP: Remember to remove a entry from the UTXOSet
	// in case that it was fully spent
	for _, tx := range transactions {
		// check if the current transaction is coinbase
		if !tx.IsCoinbase() {
			// if it isn't then range through the inputs
			for _, in := range tx.Vin {
				// delete the output of the previous transaction
				delete(u[Bytes2Hex(in.Txid)], in.OutIdx)
				// if the transaction is empty remove it entirely from the map
				if len(u[Bytes2Hex(in.Txid)]) == 0 {
					delete(u, Bytes2Hex(in.Txid))
				}
			}
			// add the outputs to the map
			u[Bytes2Hex(tx.ID)] = make(map[int]TXOutput)
			for i, out := range tx.Vout {
				// add each new output to the map
				u[Bytes2Hex(tx.ID)][i] = out
			}
		} else {
			// adds the coinbase transactions
			for i, j := range tx.Vout {
				u[Bytes2Hex(tx.ID)] = map[int]TXOutput{i: j}
			}
		}
	}
}

func (u UTXOSet) getBalance(unlockingData string) int {
	// modified version of find spendable outputs that just tells us money
	var accumulated int
	for _, txouts := range u {
		// iterates over the txouts of the current utxo
		for _, txout := range txouts {
			if txout.CanBeUnlockedWith(unlockingData) {
				// if it can be unlocked then add it to the money and add the transaction id to the map
				accumulated += txout.Value
			}
		}
	}
	return accumulated
}
