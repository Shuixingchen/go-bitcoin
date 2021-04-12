package test

import (
	"fmt"
	"letcode/programMod"
	"testing"
)

func TestPiple(t *testing.T) {
	var nums = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for n := range programMod.Pipeline(nums, programMod.Echo, programMod.Odd, programMod.Sq, programMod.Sum) {
		fmt.Println(n)
	}
}
