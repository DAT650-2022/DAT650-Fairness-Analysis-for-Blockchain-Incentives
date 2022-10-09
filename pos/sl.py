# SL pos
# time = basetime·Hash(pk, . . . )/stake
# basetime is pre-determined
# the block with smallest time will be accept.
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
  
class Blockchain:
    def __init__(self, genesis_data, basetime):
        self.chain = []
        self.chain.append(Block(genesis_data))
        self.basetime = basetime
        self.size = 0
        self.totalStake = 0
        self.winner=None
        self.lastTime=None

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
        
    # winner(miner for valid block) is the one who had the samllest time.
    def compete(self,hashStr,creator):
        time=self.basetime*int(hashStr[0:15],2)/creator.stake
        if self.lastTime==None:
            self.lastTime=time
            self.winner=creator
        else:
            if self.lastTime>time:
                self.lastTime=time
                self.winner=creator
            
    # check the block miner
    def isWinner(self,creator):
        if self.winner==creator:
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
    # candidate block
    self.candidate=None
  
    if self.blockchain != None:
      self.blockchain.totalStake += self.stake
      self.lastBlock = blockchain.lastBlock()

  def updateLast(self):
    latest = self.blockchain.lastBlock()
    if latest.height > self.lastBlock.height:
        self.lastBlock = latest

  def mine(self,seconds):
    newBlock = Block(str(self.blockchain.size), self, self.lastBlock, seconds)
    h = newBlock.pos_hash()
    self.candidate=newBlock
    self.blockchain.compete(h,self)

    
  def PoSSolver(self):
    if self.blockchain.isWinner(self):
      self.blockchain.add(self.candidate)
      self.lastBlock = self.candidate
      # stake power add 10 every time mined a new block
      self.stake=self.stake+10

def simulation(miners, number, blockchain):
    start_time = time.time()
    while blockchain.size<number:
        seconds = (time.time() - start_time)
        for miner in miners:
            miner.updateLast()
            miner.mine(seconds)
            miner.PoSSolver()

def runSimulation(times,rounds,i_miners,i_bc):
  initStakes=[0,0,0,0]
  resultStakes=[0,0,0,0]
  blockNum=[0,0,0,0]
  table=[["miner","initial stake","average final stake","average block founded"]]
  average_bn=[0,0,0,0]
  average_stakes=[0,0,0,0]

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
              initStakes[x]=miner.initialstake
              resultStakes[x]=miner.initialstake+miner.stake
              blockNum[x]= blockNum[x]+bc.checkMiner(miner)
          else:
              resultStakes[x]= resultStakes[x]+miner.stake    
              blockNum[x]= blockNum[x]+bc.checkMiner(miner)
      
  for i in range(0,4):
      item=["m"+str(i+1), initStakes[i], resultStakes[i]/times, blockNum[i]/times]
      average_bn[i]= blockNum[i]/times
      average_stakes[i]= resultStakes[i]/times
      table.append(item)
  result={
    "initial_stake": initStakes,
    "average_bn": average_bn,
    "average_stakes": average_stakes,
    "table": table
  }
  return result