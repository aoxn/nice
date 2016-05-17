package main

import (
	//"github.com/spacexnice/nice/pkg/base"
	"github.com/spacexnice/nice/pkg/algorithm"
	//"fmt"
	//"github.com/golang/glog"
	"github.com/spacexnice/nice/pkg/util"
	"os"
)

type A struct {
	S string
}

type B struct {
	M [33]A
}

func main(){
	//
	//a := int64(2 * 2.3 + 3)
	//
	//fmt.Println(a)
	//
	//bkt := base.NewBucket(false)
	//
	//x:= algorithm.NewRelateNicer(bkt).Haha()
	//for _,vx := range x{
	//
	//	glog.Infof("MERGE: %+v\n",vx)
	//}


	s,_ := os.Getwd()
	db := util.OpenInit(s)
	db.AutoMigrate(&algorithm.Record{})
	w := algorithm.NewWorker(db)
	w.FillDatabaseTest()


}