package main

import "fmt"

type A struct {
	S string
}

type B struct {
	M [33]A
}

func main(){
	b := B{}
	b.M[1].S = "3"
	a := b
	a.M[1].S = "2"
	var rt []int
	rt = append(rt,2)

	fmt.Println(b,"::::::::::",a)
}