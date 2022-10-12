# DAT650 Group Project 8 -  Do the Rich Get Richer? Fairness Analysis for Blockchain Incentives

## Report
Cant be found yet :)))

## Go Code
The Go code can be run by going into the `pow` directory and running `./pow` from the terminal. The code is also explained in the report and in comments in the code itself.

### Using simulation
Once the simulation is running by default 4 addresses are present. `a`, `b` `c` and `d`, they all have mining power 1 by default. Their mining power can be changed using `create-address [address] [power]`. This can also be used to add more miners to the mining game.

- `create-address [address] [power]` can either change the mining power or create new addresses for the mining game.
- `mine-block [address]` will mine a single block for the given address.
- `mine-blocks-pow [nr of blocks]` will mine the given number of blocks where all the addresses compete.
- `reset-bc [address]` resets the blockchain and gives the genesis block to the given address.
- `address-list` prints all the current addresses, their mining power and their current stake.
- `print-chain` prints all the blocks in the blockchain. 
- `multi-mine-blocks-pow [nr of blocks] [times to run]` will mine the given number of blocks the given number of times and print out the average, minimum and maximum stake acquired by addresses in the game. **Note:** this function does not take the current blockchain into account or commit to it when its finished.
- `help` shows all these commands in the cli.

### The new functions:
All of these functions have comments in the code itself to explain what each part does.

- `MineBlockCompete()` in `blockchain.go` is a modified version of `MineBlock()`. The difference being that this new function doesn't take in a list of transactions because thats not relevant for this simulation. Instead it takes in a list of addresses (miners) and their mining power. It then ranges through the addresses and lets them try to mine a block.
- `MineCompete()` in `block.go` is a modified version of `Mine()`. Difference being that it allows for the PoW run function to not return a block. This is important because this is how multiple miners getting their turns is implemented.
- `RunCompete()` in `proof_of_work.go` is a modified version of `Run()`. Difference being that it only attempts a single nonce for the given header. If the nonce fails it just returns `0, nil` instead of going until it finds a valid hash like `Run()` does.
- `POW_Mine_Many_Multiple()` in `mine-many.go` creates a new blockchain for each run and mines a given number of blocks. After it has finished mining the blocks it will add the balance of the miners to a map which will print out the average, max and min stake when its finished.

## Python code
### Structure
In `./pos` folder are the implementations of PoS:
- `c.py` is the implementation of Compound PoS.
- `ml.py` is the implementation of Multi-lottery PoS.
- `sl.py` is the implementation of Single-lottery PoS.
- `simulation_PoS.ipynb` is for the simulation.

More details of implementation are in the report.

### Simulation
#### Initialize
Use `simulation_PoS.ipynb` to run simulation of three types of PoS.
In the second block, the default setting is to run 1000 rounds 100 times. Set by `rounds=1000` and `times=100`.
For the miners, each miner can be set by `miner=[initial_stake, name]`, for example `m1=[10,"m1"]` which means set a miner named `m1` and got initial stake as `10`. Default setting is 4 different miners, which are:
1. Miner 1 named `m1`, staking power is `10`.
2. Miner 2 named `m2`, staking power is `5`.
3. Miner 3 named `m3`, staking power is `20`.
4. Miner 4 named `m4`, staking power is `2`.

All the miners need to be combined in one array like `miners=[m1,m2,m3,m4]` for simulation functions.
The blockchain can be set by `blockchain=[genesis data, difficulty]`, there is a default blockchain with parameters `genesis_data=0` and `difficulty=10`. <br/>
#### Run
After initializing the miners and blockchain, use different functions to run the simulations for different types of PoS. 
- `ml(...args)` for Multi-lottery PoS.
- `sl(...args)` for Single-lottery PoS.
- `c(...args)` for Compound PoS.

They all take the same arguments,
- `times` the number of running times for the simulation.
- `rounds` the number of rounds for each times of the simulation.
- `bc` the blockchain used to simulate.
- `miners` the miners. 
#### Result
The result is present in tables and plots. Three tables for each type of PoS, show the initial stake, the number of average blocks mined by each miner and their success probability of proposed a valid block. The first plot compares the average final stake from each miner, the second plot compares the number of average blocks mined by each miner, and the third one compare their success probability.