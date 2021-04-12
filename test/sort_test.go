package test

import (
	"letcode/util"
	"testing"
)

func TestBobbleSort(t *testing.T){
	input := []int{1,6,3,2,7}
	output := util.BobbleSort(input)
	t.Log(output)
}
