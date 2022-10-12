package main

// TXOutput represents a transaction output
type TXOutput struct {
	Value        int    // The transaction value
	ScriptPubKey string // The conditions to claim this output
}

// CanBeUnlockedWith checks if the output can be unlocked with the provided data
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	// TODO(student)
	// for part 1 this is a primitive store of the user name
	return out.ScriptPubKey == unlockingData
}
