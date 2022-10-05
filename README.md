# DAT650 Group Project 8 -  Do the Rich Get Richer? Fairness Analysis for Blockchain Incentives

## Report
Cant be found yet :)))

## Go Code
The Go code can be run be going into the `simulation` directory and running `./simulation` from the terminal. The code is explained in the report and in comments in the code it self.

### Using simulation
Once the simulation is running by default 2 addresses are present. `a` and `b`, they both have mining power 1 by default. Their mining power can be changed using `create-address [address] [power]`. This can also be used to add more miners to the mining game.

- `create-address [address] [power]` can either change the mining power or create new addresses for the mining game.
- `mine-block [address]` will mine a single block for the given address.
- `mine-blocks-pow [nr of blocks]` will mine the given number of blocks where all the addresses compete.
- `reset-bc [address]` resets the blockchain and gives the genesis block to the given address.
- `address-list` prints all the current addresses, their mining power and their current stake.
- `print-chain` prints all the blocks in the blockchain. 
- `multi-mine-blocks-pow [nr of blocks] [times to run]` will mine the given number of blocks the given number of times and print out the average, minimum and maximum stake acquired by addresses in the game. **Note:** this function does not take the current blockchain into account or commit to it when its finished.
- `help` shows all these commands in the cli.

### The new functions:
All of these functions have comments in the code it self to more directly explain how they work.

- `MineBlockCompete()` in `blockchain.go` is a modified version of `MineBlock()`. The difference being that this new function doesn't take in a list of transactions because thats not relevant for this simulation. Instead it takes in a list of addresses (miners) and their mining power. It then ranges through the addresses and lets them try to mine a block.
- `MineCompete()` in `block.go` is a modified version of `Mine()`. Difference being that it allows for the PoW run function to not return a block. This is important because this is how multiple miners getting their turns is implemented.
- `RunCompete()` in `proof_of_work.go` is a modified version of `Run()`. Difference being that it only attempts a single nonce for the given header. If the nonce fails it just returns `0, nil` instead of going until it finds a valid hash like `Run()` does.


## Jupyter Notebook
