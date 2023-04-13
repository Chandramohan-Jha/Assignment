package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/chandramohan/in_memory/routers"
)

func main() {
	fmt.Println("Hello")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", routers.Route).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
