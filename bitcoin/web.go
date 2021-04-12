package bitcoin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type BlockShow struct {
	ID int
	Name string
}

//显示区块链
func BlockList(ctx *gin.Context) {
	list := BC.Print()
	ctx.HTML(http.StatusOK, "blockchain.html",list)
}

//显示区块链
func BlockInfo(ctx *gin.Context) {
	id := ctx.Param("id")
	key,_ := strconv.Atoi(id)
	block := BC.FindBlock(key)
	info := make(map[string]interface{})
	info["Hash"] = block.Hash
	info["PreHash"] = block.PreHash
	info["Timestamp"] = block.Timestamp
	info["Noce"] = block.Noce
	info["High"] = BC.High()
	info["Transactions"] = block.Transactions
	ctx.HTML(http.StatusOK, "blockinfo.html",info)
}

//提交交易
func PostTransaction(ctx *gin.Context){
	var tx TxParam
	ctx.Bind(&tx)
	TxChan <- &tx

	ctx.JSON(200, gin.H{
		"code": 200,
		"message": "success",
	})
}