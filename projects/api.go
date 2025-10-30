package main

import (
	"encoding/json"
	"fmt"
)

type cars struct{

	ID int  `json:"id"`
	Color string `json:"color"`
}

func main(){

 tranf:= cars{ID:5,Color: "blue"}

  cont ,_ :=json.Marshal(tranf)

  fmt.Println(string(cont))

}