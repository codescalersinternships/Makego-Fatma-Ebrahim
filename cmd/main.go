package main

import (
	"fmt"

	makego "github.com/codescalersinternships/Makego-Fatma-Ebrahim/pkg"
)

func main() {
	err := makego.Make()
	if err != nil {
		fmt.Println(err)
	}
}
