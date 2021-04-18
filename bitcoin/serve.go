package bitcoin

import (
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"strconv"
	"sync"
)


var(
	BC = CreateBlockChain([]byte("aaa"))
	TxChan = make(chan *TxParam, 10)
	Transactions = make([]*Transaction,0)
	lockUTXO sync.RWMutex
	lockTxs sync.RWMutex
)


//客户端提交的交易信息
type TxParam struct {
	Pub string `form:"pubKey"`
	Pri string `form:"prvKey"`
	Amount string `form:"amount"`
	PayeeAddr string `form:"payeeAddr"`
}

func Serve() {
	go RunWeb()
	go AcceptTransaction()
	Mining()
}

//挖矿
func Mining() {
	for{
		lockTxs.RLock()
		txs := Transactions[:] //复制一份作打包到当前区块中
		Transactions = Transactions[0:0]//清空切片
		lockTxs.RUnlock()
		BC.AddBlock(txs)
	}
}

func AcceptTransaction() {
	for tx := range TxChan {
		amount,_ :=  strconv.Atoi(tx.Amount)
		pri,_ := hex.DecodeString(tx.Pri)
		pub,_ := hex.DecodeString(tx.Pub)
		txNew,err := NewUTXOTransaction(tx.PayeeAddr,amount, pri, pub, BC)
		if err!=nil {
			fmt.Println("valid transaction: %s", err )
		}
		lockTxs.RLock()
		Transactions = append(Transactions,txNew)
		lockTxs.RUnlock()
	}
}

//使用gin写一个web服务，接收用户转账，查看区块链信息
func RunWeb() {
	r := SetRouter()
	r.Run(":8080")
}

func SetRouter() *gin.Engine{
	r := gin.Default()
	//给模板设置自定义函数
	r.SetFuncMap(template.FuncMap{
		"ByteToString": ByteToString,
	})

	r.LoadHTMLGlob("templates/bitcoin/*")

	// 获取区块链数据路由
	r.GET("/account", Account)

	// 获取区块链数据路由
	r.GET("/blockchain", BlockList)

	// 获取某个区块信息
	r.GET("/blockinfo/:id", BlockInfo)

	//显示发起交易页面
	r.GET("/transaction", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "transaction.html", nil)
	})
	r.POST("/transaction", PostTransaction)
	return r
}