package test

import (
	"letcode/dataStruct"
	"testing"
)

/*
链表相关的算法题
*/

func TestCreateChain(t *testing.T){
	s := []int{5,7,3,1}
	list := dataStruct.NewChain()
	for _,v := range s{
		list.Insert(v)
	}
	list.Print()
	new := dataStruct.ReverseList(list)
	new.Print()
}

