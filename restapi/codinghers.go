package main

import (
	"fmt"
)

func hungrytsucess( achieve int ){

	if achieve == 5 {
		fmt.Println("5th year accomplished")
		return
	}

	fmt.Println("keep trying for sucess")

	hungrytsucess(achieve+1) // person trying himself to imprive to achive something is called recusrion 
}


func life(prob string) {

	defer func() {

		if r := recover(); r != nil { // handling exception 
			fmt.Println("recoverd ")
		}

	}()

	if prob != "happy" {
		panic(prob) // throw execption 
	}
	fmt.Println("smiles")

}

func main() {

	fmt.Println("my first wife")
	defer fmt.Println("my break up")
	fmt.Println("i am single")

	life("accident")

	hungrytsucess(0)
}
