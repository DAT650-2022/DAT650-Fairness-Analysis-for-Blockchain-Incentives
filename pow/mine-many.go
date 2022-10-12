package main

import (
	"fmt"
)

func POW_Mine_Many_Multiple(addressList map[string]int, nrToMine int, nrTimesToMine int) {
	balance := make(map[string][]int)

	for i := 0; i < nrTimesToMine; i += 1 {
		bc := NewBlockchain("nobody")

		utxos := make(UTXOSet)

		for j := 0; j < nrToMine; j += 1 {
			_, err := bc.MineBlockCompete(addressList)
			if err != nil {
				fmt.Println("error: ", err)
			} else {
				//fmt.Println("block has been added to the blockchain: ", block.String())
			}
		}
		utxos = bc.FindUTXOSet()

		for address := range addressList {
			balance[address] = append(balance[address], utxos.getBalance(address))
		}
	}

	fmt.Println("current blockreward is: ", BlockReward)
	for current, list := range balance {
		min, max := findMinAndMax(list)
		fmt.Println(" avg  |  min  |  max  |  for ", current)
		fmt.Println(average(list), " , ", min, " , ", max)

	}
}
