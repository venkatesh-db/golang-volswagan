package main

import (
	"fmt"
	"io"
	"net/http"
)

// Flow of code 

// 1. call get libarray url 
// 2. error handling 
// 3. data tranformation geting resposne []byte --> transformation --> string


func main(){

baseUrl := "https://www.redbus.in/bus-tickets/tirupathi-to-bangalore?fromCityName=Tirupati&fromCityId=71756&srcCountry=undefined&fromCityType=CITY&toCityName=Bangalore&toCityId=122&destCountry=India&toCityType=CITY&onward=30-Oct-2025&doj=30-Oct-2025&ref=home"

resp,err:=http.Get(baseUrl)

if err!=nil {
	panic(err)
}
defer resp.Body.Close()
body,_ :=io.ReadAll(resp.Body)
fmt.Println(string(body))
}

