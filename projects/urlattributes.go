package main 

import (
	  "fmt"
      "net/url"
)
func main(){

	u,_:=url.Parse("https://www.redbus.in/bus-tickets/tirupathi-to-bangalore?fromCityName=Tirupati&fromCityId=71756&srcCountry=undefined&fromCityType=CITY&toCityName=Bangalore&toCityId=122&destCountry=India&toCityType=CITY&onward=30-Oct-2025&doj=30-Oct-2025&ref=home")
	fmt.Println(u.Scheme)
	fmt.Println(u.Host)
	
}
