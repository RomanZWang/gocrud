package main

//https://www.youtube.com/watch?v=jFfo23yIWac
//https://docs.aws.amazon.com/code-library/latest/ug/go_2_dynamodb_code_examples.html

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (basics TableBasics) CreateData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data Data

	decodeError := json.NewDecoder(r.Body).Decode(&data)
	if decodeError != nil {
		log.Fatalf("failed to decode data %v\n", decodeError)
	}

	createItemError := basics.CreateTestData(data)
	if createItemError != nil {
		log.Fatalf("failed to create data %v\n", createItemError)
	}
	json.NewEncoder(w).Encode(data)
}

func (basics TableBasics) ReadData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	data, readItemError := basics.ReadTestData(params["id"])
	if readItemError != nil {
		log.Fatalf("failed to read, %v\n", readItemError)
	}
	fmt.Printf("retrieved data: %v\n", data)

	json.NewEncoder(w).Encode(data)
}

func (basics TableBasics) UpdateData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data Data

	decodeError := json.NewDecoder(r.Body).Decode(&data)
	if decodeError != nil {
		log.Fatalf("failed to decode data %v\n", decodeError)
	}

	res, createItemError := basics.UpdateTestData(data)
	if createItemError != nil {
		log.Fatalf("failed to create data %v\n", createItemError)
	}
	fmt.Printf("updated data: %v\n", res)
	json.NewEncoder(w).Encode(data)
}

func (basics TableBasics) DeleteData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	err := basics.DeleteTestData(params["id"])
	if err != nil {
		log.Fatalf("failed to delete data %v\n", params["id"])
	}
}

func main() {
	databaseClient := getDDBClient()
	basics := getTableBasics(databaseClient)

	r := mux.NewRouter()
	r.HandleFunc("/data", basics.CreateData).Methods("POST")
	r.HandleFunc("/data/{id}", basics.ReadData).Methods("GET")
	r.HandleFunc("/data", basics.UpdateData).Methods("PUT")
	r.HandleFunc("/data/{id}", basics.DeleteData).Methods("DELETE")

	fmt.Printf("server started at port 8420\n")
	log.Fatal(http.ListenAndServe(":8420", r))
}
