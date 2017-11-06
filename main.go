package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"os"
)


func main() {

	//
	route := mux.NewRouter()
	route.HandleFunc("/version", VersionHandler).
		Methods("GET")
	route.HandleFunc("/exit", ExitHandler).
		Methods("GET")
	//auth flow
	route.HandleFunc("/v1/charges", AuthauthorizeHandler).
		Methods("POST").Headers("Content-Type", "application/x-www-form-urlencoded")
	//capture flow
	route.HandleFunc("/v1/charges/{id}/capture", CaptureHandler).
		Methods("POST").Headers("Content-Type", "application/x-www-form-urlencoded")
	//refund flow
	route.HandleFunc("/v1/charges/{id}/refunds", RefundsHandler).
		Methods("POST").Headers("Content-Type", "application/x-www-form-urlencoded")
	//start server
	log.Fatal(http.ListenAndServe(":8080", route))
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"version\": \"mock-1.0.0\"}")
	log.Println("authauthorize")
}


func ExitHandler(w http.ResponseWriter, r *http.Request) {
	os.Exit(3)
}
