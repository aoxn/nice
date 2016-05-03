package app

import (
    "github.com/spacexnice/nice/pkg/algrithem"
    "github.com/spacexnice/nice/pkg/base"
    "fmt"
)

func app(){
    idx := 1491
    prd := algrithem.NewPredicator(base.NewBucket(false))
    prd.PKey3(idx).NicePrint()

    fmt.Println("\n")
    prd.Show(idx)
    //bkt := base.NewBucket(false)
    //bkt.NicePrint()

}