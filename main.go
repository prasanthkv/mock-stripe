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
	route.HandleFunc("/v1/version", VersionHandler).
		Methods("GET")
	route.HandleFunc("/server/exit", ExitHandler).
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
	fmt.Fprintln(w, "{\"version\": \"mock-1.0.1\", \"build_date\":\"nov-5-2017\"}")
	log.Println("versionn")
}


func ExitHandler(w http.ResponseWriter, r *http.Request) {
	os.Exit(3)
}
