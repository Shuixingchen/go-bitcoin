package bitcoin

import (
	"encoding/hex"
	"errors"
	"fmt"
)


type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}
type ScriptSig struct {
	Address string //付款人地址
	PubKey []byte //付款人公钥
	Sig string //本次交易签名
}

type TXInput struct{
	Txid      []byte //上个交易的transactionid
	Voutkey      int //上个交易的输出下标
	ScriptSig ScriptSig
}
type TXOutput struct {
	Value        int
	ScriptPubKey string //收款人的address
}

//创建一个铸币交易,不需要输入，输出把奖励写到地址即可
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txin := TXInput{[]byte{}, -1, ScriptSig{	Address: "",PubKey:  nil,Sig:"",}}
	txout := TXOutput{bonus, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()
	return &tx
}

/*
1.验证公钥-地址是否匹配
2.验证签名是否匹配
3.验证input是否在本地的UTXO中，防止double spending
*/
func (in *TXInput)Validate(bc *BlockChain) (int, error){
	addr := PubKeyToAddress(in.ScriptSig.PubKey)
	if addr != in.ScriptSig.Address {
		return 0,errors.New("pubkey wrong")
	}
	if RsaVerySignWithSha256(in.Txid, in.ScriptSig.Sig, in.ScriptSig.PubKey) == false {
		return 0,errors.New("sig validate fail")
	}
	out,err := bc.UTXO.FindOutput(in);
	if err!=nil {
		return 0,errors.New(err.Error())
	}
	return out.Value,nil
}


//解锁这个输出，只有解锁成功，这个输出才能使，暂时用地址来解锁。传进来的地址和输出保存的地址一致，代表成功。
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

//发起一个交易需要哪些信息
func NewUTXOTransaction(to string, amount int, pubKey []byte, prvKey []byte, bc *BlockChain) (*Transaction, error) {
	var outputs []TXOutput
	var inputs []TXInput
	from := PubKeyToAddress(pubKey)
	bc.Mux.RLock()
	acc, validOutputs := bc.UTXO.FindSpendableOutputs(from, amount)
	bc.Mux.RUnlock()

	if acc < amount {
		return nil, errors.New("ERROR: Not enough funds from "+from)
	}
	// Build a list of inputs
	for txid, outmap := range validOutputs {
		txID, _ := hex.DecodeString(txid)

		for outkey, _ := range outmap {
			//对这个input进行签名，确定是from这个人操作的
			scriptsig := ScriptSig{
				Address: from,
				PubKey:  pubKey,
				Sig:     RsaSignWithSha256(txID, prvKey),
			}
			input := TXInput{txID, outkey, scriptsig}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from}) // a change
	}

	tx := &Transaction{nil, inputs, outputs}
	tx.SetID()
	bc.UTXO.Update(tx)
	return tx, nil
}

/*
对外界接收到的交易进行验证，验证成功就可以打包到区块
1.验证每个输入是否合法
2.验证输入总额是否等于输出
 */
func (tx *Transaction) Validate(bc *BlockChain) error{
	inamount,outamount := 0,0
	for _, in := range tx.Vin {
		amount,err := in.Validate(bc)
		if err!= nil {
			return errors.New(err.Error())
		}
		inamount+=amount
	}
	for _,out := range tx.Vout {
		outamount += out.Value
	}
	if inamount != outamount {
		return errors.New("inamount != outamount")
	}
	bc.UTXO.Update(tx)
	return nil
}

func (tx *Transaction) SetID() {
	tx.ID,_ = hex.DecodeString(GetSHA256HashCode(tx))
}

func (tx *Transaction) IsCoinbase() bool{
	if len(tx.Vin) == 0 {
		return true
	} else{
		return false
	}
}