package main

import (
	"github.com/spacexnice/nice/pkg/base"
	"github.com/golang/glog"
)

var INIT_KEY = "Level:1/Index:1/Pattern:[33]/Start:1/End:33/AreaLen:33"

func main(){

	bkt    := base.NewBucket(false,-1)
	//result := bkt.NiceDebug(-1,false)
	for i := 100;i<700;i++{
		result := bkt.Nice(i)
		glog.Infoln("Result: ",i,"  ",result.Search(bkt.Balls[i]))
	}
	bkt.Statistic()

	glog.Infoln("++++++++++++++++++++++++ [RESULT] ++++++++++++++++++++++++++\nAfter PdtGroup Merge: ")
	//result.NicePrint()
}





