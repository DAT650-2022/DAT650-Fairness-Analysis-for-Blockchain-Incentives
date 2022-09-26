package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// lmao
func Execute(command string, blockchain *Blockchain, transactions []*Transaction, utxos *UTXOSet, addressList map[string][]byte) (*Blockchain, []*Transaction, *UTXOSet, map[string][]byte) {
	command = strings.TrimSuffix(command, "\n")
	args := strings.Split(command, " ")

	utxosSet := *utxos

	switch args[0] {
	case "exit":
		os.Exit(0)
	case "create-address":
		if len(args) >= 2 {
			_, pub := newKeyPair()
			addressList[args[1]] = pub
			fmt.Println("created new public address: ", args[1])
			fmt.Println("With public key: ", GetStringAddress(pub))
		} else {
			fmt.Println("not enough args")
		}
	case "print-address-list":
		for i, j := range addressList {
			fmt.Println(i, GetStringAddress(j), utxos.getBalance(HashPubKey(j)))
		}
	case "create-blockchain":
		if len(args) >= 2 {
			pubkey := cli_findKey(args[1], addressList)
			if pubkey == nil {
				fmt.Println("address not found")
				break
			}
			fmt.Println("creating a blockchain with genesis block for address", args[1])
			blockchain, err := NewBlockchain(GetStringAddress(pubkey))
			if err != nil {
				fmt.Println(err)
			}
			blockchain.FindUTXOSet()
			utxosSet = blockchain.FindUTXOSet()
			return blockchain, transactions, &utxosSet, addressList
		} else {
			fmt.Println("not enough args")
		}
	case "mine-block":
		key := cli_findKey(args[1], addressList)
		blockrw, _ := NewCoinbaseTX(GetStringAddress(key), "")
		transactions = append(transactions, blockrw)
		block, err := blockchain.MineBlock(transactions)
		utxosSet = blockchain.FindUTXOSet()
		transactions = []*Transaction{}
		if err != nil {
			fmt.Println("error: ", err)
		} else {
			fmt.Println("block has been added to the blockchain: ", block.String())
		}
		return blockchain, transactions, &utxosSet, addressList
	case "print-chain":
		fmt.Println("All blocks from the blockchain: \n", blockchain.String())
	case "help":
		fmt.Println(help())
	default:
		fmt.Println("unknown command")
	}
	return blockchain, transactions, &utxosSet, addressList
}

func cli_hashaddr(addr string) []byte {
	addByt := []byte(addr)
	pubKeyHash := Base58Decode(addByt)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	return pubKeyHash
}

func cli_findKey(addr string, addressList map[string][]byte) []byte {
	key, found := addressList[addr]
	if found {
		return key
	}
	return nil
}

func help() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("'create-address [address]': creates an address with public and private keys"))
	lines = append(lines, fmt.Sprintf("'print-address-list': lists all the addresses and their balance"))

	lines = append(lines, fmt.Sprintf("'create-blockchain [address]': creates-blockchain with genesis block to given address"))
	lines = append(lines, fmt.Sprintf("'add-transaction [from] [to] [amount]': add-transaction"))
	lines = append(lines, fmt.Sprintf("'mine-block [address]': mines a block committing the transactions and rewarding the address"))

	lines = append(lines, fmt.Sprintf("'print-chain': prints all blocks in the blockchain"))

	lines = append(lines, fmt.Sprintf("'exit: exit"))

	return strings.Join(lines, "\n")
}

func main() {
	fmt.Println(help())

	reader := bufio.NewReader(os.Stdin)

	blockchain := &Blockchain{}
	transactions := []*Transaction{}

	makeutxos := make(UTXOSet)
	utxos := &makeutxos

	addressList := make(map[string][]byte)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		input = strings.TrimSuffix(input, "\n")
		if input == "" {
			continue
		}

		blockchain, transactions, utxos, addressList = Execute(input, blockchain, transactions, utxos, addressList)
		fmt.Println("")
	}
}
