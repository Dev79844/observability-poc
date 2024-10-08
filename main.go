package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
)

func main(){
	router := mux.NewRouter()

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	log.Println("server started on port 3000")
	err := http.ListenAndServe(":3000", router)
	log.Fatal(err)
}