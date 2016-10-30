package main

import (
	"github.com/spacexnice/nice/pkg/api"
	"fmt"
	"sort"
)

var INIT_KEY = "Level:1/Index:1/Pattern:[33]/Start:1/End:33/AreaLen:33"

func main(){

	//bkt    := base.NewBucket(false,-1)
	////result := bkt.NiceDebug(-1,false)
	//for i := 100;i<700;i++{
	//	result := bkt.Nice(i)
	//	glog.Infoln("Result: ",i,"  ",result.Search(bkt.Balls[i]))
	//}
	//bkt.Statistic()
	//
	//glog.Infoln("++++++++++++++++++++++++ [RESULT] ++++++++++++++++++++++++++\nAfter PdtGroup Merge: ")
	////result.NicePrint()
	bkt := api.LoadBucket(1,false)
	api.Pick(bkt)
}

func area(bkt *api.Bucket){
	e1 := bkt.Estimate(&api.K3Policy{Target:"2:2:2"})
	e6 := bkt.Estimate(&api.K3Policy{Target:"2:1:3"})
	e7 := bkt.Estimate(&api.K3Policy{Target:"2:3:1"})
	e2 := bkt.Estimate(&api.K3Policy{Target:"1:2:3"})
	e3 := bkt.Estimate(&api.K3Policy{Target:"1:3:2"})
	e4 := bkt.Estimate(&api.K3Policy{Target:"3:2:1"})
	e5 := bkt.Estimate(&api.K3Policy{Target:"3:1:2"})
	e8 := bkt.Estimate(&api.K3Policy{Target:"3:3:0"})
	k3 := api.Estimators{e1,e2,e3,e4,e5,e6,e7,e8}
	sort.Sort(k3)
	fmt.Println(k3)
}





