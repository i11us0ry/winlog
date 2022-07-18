package main

import (
	"flag"
	"winlog/src"
)

func main(){
	PathPtr := flag.String("p", "", "mimi path")
	flag.Parse()
	src.Start(*PathPtr)
}