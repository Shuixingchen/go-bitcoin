package bitcoin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
)

type BlockShow struct {
	ID   int
	Name string
}

// 显示区块链
func BlockList(ctx *gin.Context) {
	list := BC.Print()
	ctx.HTML(http.StatusOK, "blockchain.html", list)
}

// 显示区块链
func BlockInfo(ctx *gin.Context) {
	id := ctx.Param("id")
	key, _ := strconv.Atoi(id)
	block := BC.FindBlock(key)
	info := make(map[string]interface{})
	info["Hash"] = block.Hash
	info["PreHash"] = block.PreHash
	info["Timestamp"] = block.Timestamp
	info["Noce"] = block.Noce
	info["High"] = BC.High()
	info["Transactions"] = block.Transactions
	ctx.HTML(http.StatusOK, "blockinfo.html", info)
}

// 提交交易
func PostTransaction(ctx *gin.Context) {
	var tx TxParam
	err := ctx.Bind(&tx)
	if err != nil {
		logs.WithFields(logs.Fields{"method": "PostTransaction"}).Error(err)
	}
	TxChan <- &tx

	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}

// 创建新账号
func Account(ctx *gin.Context) {
	pri, pub := GenRsaKey()
	ctx.HTML(http.StatusOK, "account.html", map[string][]byte{"pri": pri, "pub": pub})
}
