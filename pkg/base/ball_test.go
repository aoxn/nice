package base

import (
	"testing"
	"github.com/golang/glog"
)

func TestKeyPartition(t *testing.T) {
	bkt := NewBucket(false, -1)
	for i := 100; i < 1900; i++ {
		result := bkt.Nice(i)
		t.Log("Result: ", i, "  ", result.Search(bkt.Balls[i]))
	}

	bkt.Statistic()

	glog.Infoln("++++++++++++++++++++++++ [RESULT] ++++++++++++++++++++++++++\nAfter PdtGroup Merge: ")
	//result.NicePrint()
}
