package main

import "fmt"

// wife --> Crud
// crud --> create read update delete

//  venkatesh --> seetha
//        ||
//  common database --> mysql postrgess

// scenario venkatesh
// we buy 1 cr @indiranagar -thinking
// wife --> create insert in to databse --> Husband
// wife --> update --> 1.5 cr @jayangar --> Wife
// parent husband and wife              --> Read
// authentication ---> builder
// security ---> money security
// fight --> property site.             --> delete
// encoding --> msg --> wife scolding encoding decoding her mother


//Schema human life

// Mysql                    postgressdb                    redis 

// common wife-husband       specifc habbits husband.       repeatly husband wife is fighting


//   Kafka -->  Work ,  entirament , social media 
//  human -->  consumer  consumer      consumer

//                 wife    parents     friend's --> producing

//  husband -->  consumer  consumer    consumer  
// husband --> recive msg from  wife    parents     friend's  
// husband need to be multiple  goroutines 


var conversation map[string][]string =make(map[string][]string)


func createhudbandwifeconvesation( propery string , con []string){
  conversation[propery]=con
  fmt.Println("insert happened", conversation[propery])
}

func readhudbandwifeconvesation(con string) []string {
	 fmt.Println("read happened")

	return conversation[con]
}

func updatehudbandwifeconvesation(prov string,con string){

	fmt.Println("update happened")
	conversation[prov]=append(  conversation[prov],con)
	fmt.Println("update happened",conversation[prov])
}

func deletehudbandwifeconvesation(prop string){

	 delete( conversation,prop)
	 fmt.Println("after deleting",conversation)
}


func main(){

   createhudbandwifeconvesation("property",[]string{"did u check legal","verifieddocuments"})
   ret:=readhudbandwifeconvesation("property")
   fmt.Println("reading done",ret)

   updatehudbandwifeconvesation("property","can it in 12 months")
   deletehudbandwifeconvesation("property")

}

