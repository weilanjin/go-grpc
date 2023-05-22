package util

import (
	"math/rand"
	"time"
)

// RandNumber 指定范围的随机数
func RandNumber(start, end int) func() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return func() int {
		if start >= end {
			return end
		}
		num := r.Intn(end-start) + start
		return num
	}
}
