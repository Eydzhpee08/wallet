package main

import (
	"log"
	// "os"
	// "github.com/google/uuid"
	"fmt"
	// "github.com/Eydzhpee08/wallet/pkg/types"
	"github.com/Eydzhpee08/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	// _, err := svc.RegisterAccount("+992929582003")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// _, err = svc.RegisterAccount("+992929580003")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// _, err = svc.RegisterAccount("+992907107797")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	path := "../data/accounts.txt"

	err := svc.ImportFromFile(path)

	if err != nil {
		log.Println(err)
	}
	fmt.Println("Success")
	// err = svc.ExportToFile(path)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println("Success")

	// err = svc.Deposit(account.ID, 100)
	// if err != nil {
	// switch err {
	// case wallet.ErrAmountMustBePositive:
	// fmt.Println(wallet.ErrAmountMustBePositive)
	// case wallet.ErrAccountNotFound:
	// fmt.Println(wallet.ErrAccountNotFound)
	// }
	// return
	// }
	// fmt.Println(account.Balance)
	//
	// payment, err := svc.Pay(account.ID, 50, types.PaymentCategoryFun)
	// if err != nil {
	// switch err {
	// case wallet.ErrAmountMustBePositive:
	// fmt.Println(wallet.ErrAmountMustBePositive)
	// case wallet.ErrNotEnoughBalance:
	// fmt.Println(wallet.ErrNotEnoughBalance)
	// case wallet.ErrAccountNotFound:
	// fmt.Println(wallet.ErrAccountNotFound)
	// }
	// return
	// }
	// fmt.Println(payment)

	// file, err := os.Open("../data/readme.txt")
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// defer func() {
	// 	err := file.Close()
	// 	if err != nil {
	// 		log.Print(err)
	// 	}
	// }()

	// buf := make([]byte, 4)
	// read, err := file.Read(buf)
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// data := string(buf[:read])
	// log.Print(data)
}
