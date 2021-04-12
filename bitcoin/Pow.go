package bitcoin

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
)

type Pow struct {
	block *Block
	target *big.Int
	targetBits int
}

//传入区块头,目标值
func NewPoW(b *Block,targetBits int) *Pow {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits)) //bigInt左移动
	pow := &Pow{b,target,targetBits}
	return pow
}

func (p *Pow)PrepareData(nonce int) []byte{
	data := bytes.Join(
		[][]byte{
			p.block.PreHash,
			p.block.HashTransactions(),
			IntToHex(p.block.Timestamp),
			IntToHex(int64(p.targetBits)),
			IntToHex(int64(nonce)),
		},[]byte{})
	return data
}

//挖矿，返回noce和hash
func (p *Pow)Run()(int,[]byte){
	var hashInt big.Int
	var hash [32]byte
	noce := 0
	for noce < math.MaxInt64 {
		data := p.PrepareData(noce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(p.target) == -1 {
			break
		}else{
			noce++
		}
	}
	return noce,hash[:]
}