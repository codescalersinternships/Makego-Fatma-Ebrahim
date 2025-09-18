package main

import (
	"fmt"

	makego "github.com/codescalersinternships/Makego-Fatma-Ebrahim/pkg"
)

func main() {
	targets, err := makego.Parse()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v \n",targets)
}