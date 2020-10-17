package main

import (
	"log"
	"fmt"
	"github.com/Eydzhpee08/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}

	path := "../data/accounts.txt"

	err := svc.ImportFromFile(path)

	if err != nil {
		log.Println(err)
	}
	fmt.Println("Success")

}
