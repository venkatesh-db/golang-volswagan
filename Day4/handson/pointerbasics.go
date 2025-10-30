package main

import "fmt"

type batman struct {
	greatcoder string
	expcoding  int
}

func house() {

	// new keyword to allocate single memory in heap but its one time allocation

	var location *string = new(string) // one time storage in heap memory

	*location = "railway station --> first left walk 1 km -> right walk 500m -> blue house with red roof"

	fmt.Println("House location is:", *location)

	var beating *int = new(int)

	*beating = 100

	fmt.Println("teacher heart is beating at:", *beating)

	// newkeywod with continous memory allocation

	var friends *[50]string = new([50]string)

	friends[0] = "venkatesh"
	friends[1] = "james"
	friends[2] = "bond"

	fmt.Println("My friends are:", friends[0], friends[1], friends[2])
}

func newvsmake() {

	// 50 friends can we acess only single pointer
	// we cannot append it
	var friends *[50]string = new([50]string)

	friends[0] = "venkatesh"
	friends[1] = "james"
	friends[2] = "bond"

	// 50 friends can we acess only single slice
	// make we can append it 
	var friendss []string = make([]string, 50)

	friendss[0] = "venkatesh"
	friendss[1] = "james"
	friendss[2] = "bond"

	friendss = append(friendss, "ravi kanth")
	fmt.Println("My friends are using new:", friends[0], friends[1], friends[2])
	fmt.Println("My friends are using make:", friendss[0], friendss[1], friendss[2])
}

func main() {

	// Pointer basics
	var smiles int = 5
	var jamesbond *int = &smiles

	println("Value of smiles:", smiles)
	println("Address of smiles:", &smiles)
	println("Value of jamesbond (address of smiles):", jamesbond)
	println("Value pointed to by jamesbond:", *jamesbond)

	var venkatesh *batman = &batman{
		greatcoder: "venkatesh",
		expcoding:  10,
	}

	println("Value of venkatesh:", venkatesh, *&venkatesh.expcoding)

	house()
	newvsmake()
}
