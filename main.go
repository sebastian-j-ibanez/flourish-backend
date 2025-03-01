package main

import (
	"fmt"

	"github.com/sebastian-j-ibanez/flourish-backend/code"
)

func main() {
	for range 10 {
		fmt.Println(code.GenerateCode())
	}
}
