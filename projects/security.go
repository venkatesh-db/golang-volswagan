package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func main() {

	token := make([]byte, 16)

	_, err := rand.Read(token)

	if err != nil {
		panic(err)
	}

	fmt.Println("secure token", hex.EncodeToString(token))

}
