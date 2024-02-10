package main

import (
	"contagio/contagio/cnc/utils"
	"fmt"
	"os"
)

func main() {
	if os.Args[1] == "hash" {
		fmt.Print(utils.Sha3(os.Args[2])) // login
	}
}
