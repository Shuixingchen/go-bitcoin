package test

import (
	"bitcoin/bitcoin"
	"bitcoin/db"
	"fmt"
	"testing"
)

func TestDb(t *testing.T) {
	mydb := db.GetDB()
	defer mydb.Close()
	mydb.Put([]byte("key"), []byte("valueyyy"))
	data := mydb.Get([]byte("key"))
	fmt.Println(string(data))
}

func TestSerialize(t *testing.T) {
	b := bitcoin.Block{
		Hash:         nil,
		PreHash:      nil,
		Transactions: nil,
		Timestamp:    55555,
		Noce:         555,
	}
	data := b.Serialize()
	res := bitcoin.DeserializeBlock(data)
	fmt.Print(res)
}
