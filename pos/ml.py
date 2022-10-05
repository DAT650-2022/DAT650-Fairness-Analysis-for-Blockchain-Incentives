# multi-lottery
import hashlib as hasher
import random
import time
def hashbits(input):
    hash_obj = hasher.sha256()
    inputbytes = input.encode()
    #print(type(inputbytes))
    hash_obj.update(inputbytes)
    hashbytes = hash_obj.digest()
    return ''.join(f'{x:08b}' for x in hashbytes)

def hash(input):
    hash_obj = hasher.sha256()
    inputbytes = input.encode()
    #print(type(inputbytes))
    hash_obj.update(inputbytes)
    return hash_obj.hexdigest()

class Block:
    def __init__(self, data, creator=None, previous=None, time=0):
        self.data = data
        if previous is None:
            self.previous = None
            self.previous_hash = ""
            self.creator = Minter(0 , "0")
            self.height = 0
        else:
            self.previous = previous
            self.previous_hash = previous.hash
            self.creator = creator
            self.height = previous.height+1
        self.timestamp = time
        self.hash = self.hash_block()
        self.children = []

    def pos_hash(self):
        return hashbits(self.creator.name + self.previous_hash + str(self.timestamp))

    def hash_block(self):
        return hashbits(self.creator.name + str(self.data) + self.previous_hash + str(self.timestamp))

    def print(self):
      print(self.data + " "+ self.creator.name + " " + str(self.height))
        
class Blockchain:
    def __init__(self, genesis_data, difficulty):
        self.chain = []
        self.chain.append(Block(genesis_data))
        self.difficulty = difficulty
        self.size = 0
        self.totalStake = 0

    def lastBlock(self):
      max = self.chain[0].height
      for block in self.chain:
        if block.height > max:
          max = block.height
      maxes = [block for block in self.chain if block.height == max]
      r = random.choices(maxes, k=1)
      return r[0]

    def lastBlocks(self):
      max = self.chain[0].height
      for block in self.chain:
        if block.height > max:
          max = block.height
      maxes = [block for block in self.chain if block.height == max]
      return maxes
        
    def add(self, newBlock):
        self.chain.append(newBlock)
        newBlock.previous.children.append(newBlock)
        self.size +=1
        #newBlock.creator.stake+=1
    
    def isSmaller(self, hashStr, creator):
      #add this function
      # use int(hashStr[0:15],2) to convert the first 15 bits to int 
      # compare it with the difficulty, multiplicated by the creators stake
      if int(hashStr[0:15],2) < self.difficulty * (creator.stake+self.checkMiner(creator)):
        return True
      return False

    
    def checkMiner(self, miner, last=None):
      if last == None:
        last = self.lastBlock()
      count = 0
      while last!=None:
        if last.creator == miner:
          count += 1
        last = last.previous
      return count

class Minter:
  def __init__(self, stake, name, blockchain=None):
    self.initialstake=stake
    self.stake = stake
    self.name = name
    self.blockchain = blockchain
    
    if self.blockchain != None:
      self.blockchain.totalStake += self.stake
      self.lastBlock = blockchain.lastBlock()

  def updateLast(self):
    latest = self.blockchain.lastBlock()
    if latest.height > self.lastBlock.height:
        self.lastBlock = latest

  def PoSSolver(self, seconds):
    newBlock = Block(str(self.blockchain.size), self, self.lastBlock, seconds)
    h = newBlock.pos_hash()
    if self.blockchain.isSmaller(h,self):
      self.blockchain.add(newBlock)
      self.lastBlock = newBlock
      # stake power add 10 every time mined a new block
      self.stake=self.stake+10


# simulation
def simulation(miners, number, blockchain):
    start_time = time.time()
    while blockchain.size<number:
        for miner in miners:
            seconds = (time.time() - start_time)
            miner.updateLast()
            miner.PoSSolver(seconds)

def runSimulation(times,rounds,i_miners,i_bc):
  ml_initStakes=[0,0,0,0]
  ml_resultStakes=[0,0,0,0]
  ml_blockNum=[0,0,0,0]
  ml_table=[["miner","initial stake","average final stake","average block founded"]]
  ml_average_bn=[0,0,0,0]
  ml_average_stakes=[0,0,0,0]

  bbc=i_bc.copy()
  mminers=i_miners.copy()

  for i in range(0,times):
      bc=Blockchain(bbc[0] , bbc[1])
      miners=[]
      for m in mminers:
            newMiner=Minter(m[0],m[1],bc)
            miners.append(newMiner)
      simulation(miners,rounds,bc)
      for x, miner in enumerate(miners):
          if i==1: 
              ml_initStakes[x]=miner.initialstake
              ml_resultStakes[x]=miner.initialstake+miner.stake
              ml_blockNum[x]=ml_blockNum[x]+bc.checkMiner(miner)
          else:
              ml_resultStakes[x]=ml_resultStakes[x]+miner.stake    
              ml_blockNum[x]=ml_blockNum[x]+bc.checkMiner(miner)
      print("round: ",i)
      
  for i in range(0,4):
      name="m"+str(i+1)
      item=[name,ml_initStakes[i],ml_resultStakes[i]/times,ml_blockNum[i]/times]
      ml_average_bn[i]=ml_blockNum[i]/times
      ml_average_stakes[i]=ml_resultStakes[i]/times
      ml_table.append(item)
  result={
    "initial_stake":ml_initStakes,
    "average_bn":ml_average_bn,
    "average_stakes":ml_average_stakes,
    "table":ml_table
  }
  return result