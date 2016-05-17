package app

import (
    "github.com/spacexnice/nice/pkg/algorithm"
    "github.com/spacexnice/nice/pkg/base"
    "fmt"
)

func app(){
    idx := 1491
    prd := algorithm.NewPartitionNicer(base.NewBucket(false,0))
    prd.PKey3(idx).NicePrint()

    fmt.Println("\n")
    prd.Show(idx)
    //bkt := base.NewBucket(false)
    //bkt.NicePrint()

}