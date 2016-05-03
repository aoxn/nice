package main

import (

    "runtime"
    "github.com/spacexnice/nice/cmd/roller/app"
)

func main(){
    runtime.GOMAXPROCS(runtime.NumCPU())
    s := app.NewNiceServer()
    s.AddFlags()
    s.Run()
}

