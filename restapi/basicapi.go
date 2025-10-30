package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// to create api we need net/http Libraries
// inside api code flow --> set of libraries
// each librariy we need to error handling

/*

http://localhost:9090/user

curl -X GET http://localhost:9090/user \
-H "Content-Type: application/json" \
-d '{"name":"Venkatesh","email":"venky@gmail.com"}'

{"email":"venky@gmail.com","name":"Venkatesh","status":"sucess"}

*/

func primecar(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(res, `{"message":"hi bmw car"}`)

}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email`
}

func postHandler(res http.ResponseWriter, req *http.Request) {

	var smiles User
	err := json.NewDecoder(req.Body).Decode(&smiles)

	if err != nil {
		http.Error(res, err.Error(), http.StatusMethodNotAllowed)
	}
	response := map[string]string{
		"status": "sucess",
		"name":   smiles.Name,
		"email":  smiles.Email,
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(response)
}

func main() {

	http.HandleFunc("/hello", primecar)
	http.HandleFunc("/user", postHandler)

	fmt.Println("server is running")
	http.ListenAndServe(":9090", nil)
}
