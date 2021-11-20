package db

import (
	"log"

	"github.com/boltdb/bolt"
)

var (
	Table = []byte("MyBlocks")
)

type BlockDD struct {
	db *bolt.DB
}

func GetDB() *BlockDD {
	// path, _ := os.Getwd()
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	CreateTable(db) // 检查是否有生成MyBlocks表
	return &BlockDD{db: db}
}

func (bd *BlockDD) Close() {
	bd.db.Close()
}

func CreateTable(db *bolt.DB) {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(Table)
		if b == nil {
			_, err := tx.CreateBucket(Table)
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

// 插入或更新键值对
func (bd *BlockDD) Put(key, value []byte) {
	err := bd.db.Update(func(tx *bolt.Tx) error {
		// 取出叫"MyBucket"的表
		b := tx.Bucket([]byte("MyBlocks"))
		// 往表里面存储数据
		if b != nil {
			err := b.Put(key, value)
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (bd *BlockDD) Get(key []byte) []byte {
	var data []byte
	err := bd.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBlocks"))
		// 往表里面存储数据
		if b != nil {
			data = b.Get(key)
		}
		return nil
	})
	// 查询数据库失败
	if err != nil {
		log.Fatal(err)
	}
	return data
}
