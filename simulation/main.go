package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// lmao
func Execute(command string, blockchain *Blockchain, transactions []*Transaction, utxos *UTXOSet, addressList map[string]int) (*Blockchain, []*Transaction, *UTXOSet, map[string]int) {
	command = strings.TrimSuffix(command, "\n")
	args := strings.Split(command, " ")

	utxosSet := *utxos

	switch args[0] {
	case "exit":
		os.Exit(0)
	case "create-address":
		if len(args) >= 3 {
			number, err := strconv.Atoi(args[2])
			if err != nil {
				fmt.Println(err)
				break
			}
			addressList[args[1]] = number
			fmt.Println("created miner: ", args[1], "with mining power: ", number)
		} else {
			fmt.Println("not enough args")
		}
	case "print-address-list":
		fmt.Println("address | Mining Power | Balance")
		for i, j := range addressList {
			fmt.Println(i, j, utxos.getBalance(i))
		}
	case "create-blockchain":
		if len(args) >= 2 {
			fmt.Println("creating a blockchain with genesis block for address", args[1])
			blockchain := NewBlockchain(args[1])
			blockchain.FindUTXOSet()
			utxosSet = blockchain.FindUTXOSet()
			return blockchain, transactions, &utxosSet, addressList
		} else {
			fmt.Println("not enough args")
		}
	case "mine-blocks":
		nrToMine, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err)
			break
		}
		for i := 0; i < nrToMine; i += 1 {
			block, err := blockchain.MineBlockCompete(addressList)
			if err != nil {
				fmt.Println("error: ", err)
			} else {
				fmt.Println("block has been added to the blockchain: ", block.String())
			}
		}
		utxosSet = blockchain.FindUTXOSet()
		transactions = []*Transaction{}
		return blockchain, transactions, &utxosSet, addressList
	case "mine-block":
		// take this out later
		blockrw := NewCoinbaseTX(args[1], "")
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

	addressList := make(map[string]int)

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
