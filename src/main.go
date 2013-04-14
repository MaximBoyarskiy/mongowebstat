package main

import (
	"./mongowebstat"
)

func main() {
	if err := mongowebstat.LoadConfig(); err != nil {
		panic(err)
	}
	mongowebstat.Start()

}
