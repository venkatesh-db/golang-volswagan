package main

import (
	"fmt"
        "os"
		"bufio"
)

func readingcontent(){

	file, err := os.Open("sucess.txt")

	if err!=nil {
		panic(err)
	}

	defer file.Close()

   scanner := bufio.NewScanner(file)

   if err := scanner.Err(); err != nil {
    panic(err)
   }

   for scanner.Scan() {
       line := scanner.Text()
	   fmt.Println(line)
   }



}


func main(){

	data,err := os.ReadFile("sucess.txt")
	if err != nil {

		panic(err)
	}
	fmt.Println(string(data))

	readingcontent()

}