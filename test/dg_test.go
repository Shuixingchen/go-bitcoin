package test

import (
	"fmt"
	"letcode/bitcoin/db"
	"testing"
)

func TestCreateDb(t *testing.T) {
	db.GetConfig()
	db := db.CreateDb()
	fmt.Println(db)
}
