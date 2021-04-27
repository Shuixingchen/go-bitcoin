package test

import (
	"bitcoin/bitcoin"
	db2 "bitcoin/db"
	"fmt"
	"testing"
)

func TestDb(t *testing.T) {
	db := db2.GetDb()
	defer db.Close()
	db.Put([]byte("key"),[]byte("valueyyy"))
	data := db.Get([]byte("key"))
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
