package bitcoin

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

const(
	bonus = 10
	targetBits = 24 //难度
)


type Block struct {
	Hash []byte
	PreHash []byte
	Transactions []*Transaction
	Timestamp int64
	targetBits int
	Noce int
}

type BlockChain struct {
	blocks []*Block
	UTXO *UTXO
	Mux sync.RWMutex //读写锁
}

func CreateBlock(data []*Transaction, preHash []byte) *Block{
	newBlock := &Block{
		PreHash:preHash,
		Transactions:data,
		Timestamp:time.Now().Unix(),
	}
	//工作量证明
	pow := NewPoW(newBlock,targetBits)
	noce,hash := pow.Run()
	newBlock.Hash = hash
	newBlock.Noce = noce
	return newBlock
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

//create
func CreateBlockChain(pub []byte) *BlockChain{
	address := PubKeyToAddress(pub)
	coinBase := NewCoinbaseTX(address, "")
	genesis := CreateBlock([]*Transaction{coinBase}, []byte(""))
	bc := &BlockChain{[]*Block{genesis}, nil, sync.RWMutex{}}
	utxo := bc.FindUTXO()
	bc.UTXO = &utxo
	return bc
}

func (b *BlockChain)AddBlock(data []*Transaction){
	preBlock := b.blocks[len(b.blocks)-1]
	block := CreateBlock(data, preBlock.Hash)
	b.Mux.Lock()
	b.blocks = append(b.blocks,block)
	b.Mux.Unlock()
}

/*
通过对区块进行迭代找到所有未花费输出
transactionid=>uoutputs
 */

func (bc *BlockChain)FindUTXO() UTXO {
	bc.Mux.RLock()
	utxo := make(UTXO)
	spentTXOs := make(map[string][]int)//transactionId=>[output的下标]

	for _,block := range bc.blocks {
		for _,tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
			unsedOutput := make(map[int]TXOutput,0)
		Outputs:
			//遍历output,找出没有用的output
			for outIdx, out:= range tx.Vout {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				//this output was not used
				unsedOutput[outIdx] = out
			}
			utxo[txID] = unsedOutput

			//遍历input,找出已经使用的output
			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					//find spend output
					inTxID := hex.EncodeToString(in.Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Voutkey)
					//delete has put in utxo
					if outmap,ok := utxo[inTxID]; ok{
						for k,_:= range outmap {
							if k == in.Voutkey {
								delete(utxo[inTxID],k) //从utxo中删除这个outIdx
								if len(utxo[inTxID]) == 0 {
									delete(utxo, inTxID)
								}
							}
						}
					}
				}
			}
		}
	}
	bc.Mux.RUnlock()
	return utxo
}


func (bc *BlockChain) FindBlock(key int) *Block{
	bc.Mux.RLock()
	block := bc.blocks[key]
	bc.Mux.RUnlock()
	return block
}

func (bc *BlockChain) High() int{
	bc.Mux.RLock()
	heigh := len(bc.blocks)
	bc.Mux.RUnlock()
	return heigh
}
func (bc *BlockChain)Print() map[int]string{
	bc.Mux.RLock()
	list := make(map[int]string)
	for k,block := range bc.blocks {
		s := ""
		s += fmt.Sprintf("pre hash:%x\n",block.PreHash)
		//s += fmt.Sprintf("data:%s\n",block.HashTransactions())
		s += fmt.Sprintf("hash:%x\n",block.Hash)
		s += fmt.Sprintf("createTime:%d\n",block.Timestamp)
		s += fmt.Sprintf("\n")
		list[k] = s
	}
	bc.Mux.RUnlock()
	return list
}
