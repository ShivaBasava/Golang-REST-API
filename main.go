package main

import (
    "net/http"
    "log"
    "github.com/ShivaBasava/Golang-REST-API/handlers"    
    "github.com/ShivaBasava/Golang-REST-API/models"
    "github.com/gorilla/mux"
)


func main() {
// The router is now formed by calling the `NewRouter` constructor function
    r := mux.NewRouter()
    r.HandleFunc("/", hello).Methods("GET")

    r.HandleFunc("/api/cartoons", getCartoons).Methods("GET")
    r.HandleFunc("/api/cartoons/{id}", getCartoon).Methods("GET")
    r.HandleFunc("/api/cartoons", createCartoon).Methods("POST")
    r.HandleFunc("/api/cartoons/{id}", updateCartoon).Methods("PUT")
    r.HandleFunc("/api/cartoons/{id}", deleteCartoon).Methods("DELETE")

    log.Fatal(http.ListenAndServe(":8000", r))
}
