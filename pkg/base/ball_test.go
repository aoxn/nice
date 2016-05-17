package base

import (
	"testing"
	"fmt"
)

func TestKeyPartition(t *testing.T) {
	ball := Ball{
		Reds:[]int{2,4,6,8,10,12},
	}
	m := ball.KeyPartition("11:11:11",1)
	if m !="GRP/5:1:0"{
		t.Fail()
	}

	m = ball.KeyPartition("3:4:4",1)
	fmt.Println(m)
	if m !="GRP/1:2:2"{
		t.Fail()
	}

	m = ball.KeyPartition("3:4:4",12)
	fmt.Println(m)
	if m !="GRP/1:0:0"{
		t.Fail()
	}
}
