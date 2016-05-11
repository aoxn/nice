package main

import (

    "runtime"
    "github.com/spacexnice/nice/cmd/roller/app"
    "github.com/spacexnice/nice/pkg/base"
    "github.com/spacexnice/nice/pkg/algorithm"
    //"fmt"
)

func main(){
    //bkt := base.NewBucket(false)
    //fmt.Println(bkt.Balls[len(bkt.Balls)-1].Attr.CoRelate3)
    algorithm.NewPartitionNicer(base.NewBucket(false))
    runtime.GOMAXPROCS(runtime.NumCPU())
    s := app.NewNiceServer()
    s.AddFlags()
    s.Run()
}



