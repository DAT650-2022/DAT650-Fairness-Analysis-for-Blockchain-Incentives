package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// lmao
func Execute(command string, blockchain *Blockchain, utxos *UTXOSet, addressList map[string]int) (*Blockchain, *UTXOSet, map[string]int) {
	command = strings.TrimSuffix(command, "\n")
	args := strings.Split(command, " ")

	utxosSet := *utxos

	switch args[0] {
	case "exit":
		os.Exit(0)
	case "create-address", "change-address":
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
	case "address-list", "print-address-list", "a-l":
		fmt.Println("address | Mining Power | Balance")
		for i, j := range addressList {
			fmt.Println(i, j, utxos.getBalance(i))
		}
	case "reset-bc", "create-blockchain":
		if len(args) >= 2 {
			fmt.Println("creating a blockchain with genesis block for address", args[1])
			blockchain := NewBlockchain(args[1])
			blockchain.FindUTXOSet()
			utxosSet = blockchain.FindUTXOSet()
			return blockchain, &utxosSet, addressList
		} else {
			fmt.Println("not enough args")
		}
	case "mine-blocks-pow", "mb-pow":
		if len(args) >= 2 {
			nrToMine, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println(err)
				break
			}
			for i := 0; i < nrToMine; i += 1 {
				_, err := blockchain.MineBlockCompete(addressList)
				if err != nil {
					fmt.Println("error: ", err)
				} else {
					//fmt.Println("block has been added to the blockchain: ", block.String())
				}
			}
			utxosSet = blockchain.FindUTXOSet()
			return blockchain, &utxosSet, addressList
		} else {
			fmt.Println("not enough args")
		}
	case "multi-mine-blocks-pow", "mmb-pow":
		if len(args) >= 3 {
			nrToMine, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println(err)
				break
			}
			nrTimesToMine, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println(err)
				break
			}
			POW_Mine_Many_Multiple(addressList, nrToMine, nrTimesToMine)
		} else {
			fmt.Println("not enough args")
		}
	case "mine-block":
		// take this out later
		blockrw := NewCoinbaseTX(args[1], "")
		transactions := []*Transaction{blockrw}
		block, err := blockchain.MineBlock(transactions)
		if err != nil {
			fmt.Println("error: ", err)
		} else {
			utxos.Update(transactions)
			fmt.Println("block has been added to the blockchain: ", block.String())
		}
		return blockchain, utxos, addressList
	case "print-chain":
		fmt.Println("All blocks from the blockchain: \n", blockchain.String())
	case "help":
		fmt.Println(help())
	default:
		fmt.Println("unknown command")
	}
	return blockchain, &utxosSet, addressList
}

func help() string {
	var lines []string
	lines = append(lines, "'create-address [address] [power]': creates an address with mining power, can also be used to change mining power")
	lines = append(lines, "'address-list': lists all the addresses and their balance")

	lines = append(lines, "'reset-bc [address]': resets the blockchain and gives the genesis block to the address")
	lines = append(lines, "'mine-block [address]': mines a single block for the given address")
	lines = append(lines, "'mine-blocks-pow [nr]': PoW mines given number of blocks for all addresses")
	lines = append(lines, "'multi-mine-blocks-pow [nr] [nrT]': PoW mines [nr] of blocks for all addresses [nrT] times ")

	lines = append(lines, "'print-chain': prints all blocks in the blockchain")

	lines = append(lines, "'exit: exit")

	return strings.Join(lines, "\n")
}

func main() {
	fmt.Println(help())

	reader := bufio.NewReader(os.Stdin)

	//blockchain := &Blockchain{}

	makeutxos := make(UTXOSet)
	utxos := &makeutxos

	addressList := make(map[string]int)

	addressList["a"] = 2
	addressList["b"] = 2

	blockchain := NewBlockchain("a")

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

		blockchain, utxos, addressList = Execute(input, blockchain, utxos, addressList)
		fmt.Println("")
	}
}
