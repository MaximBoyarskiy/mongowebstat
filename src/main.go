package main

import (
	"./mongowebstat"
	"flag"
)

func main() {
	if err := mongowebstat.LoadConfig(); err != nil {
		panic(err)
	}
	httpPtr := flag.String("http", ":8080", "http listen address")
	flag.Parse()
	mongowebstat.ServerStart(httpPtr)

}
