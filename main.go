package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "os"
)

const Version = "mock-1.0.2"
const Build = "nov-6-2017"

//
// Main Method
//
func main() {
    //
    route := mux.NewRouter()
    //
    // Admin flows
    //
    route.HandleFunc("/v1/version", VersionHandler).
        Methods("GET")
    route.HandleFunc("/admin/exit", ExitHandler).
        Methods("GET")

    //
    // mocks
    //

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

//
// Handle version request
//
func VersionHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "{\"version\": \"" + Version + "\", \"build_date\":\"" + Build + "\"}")
    log.Println("VersionHandler:Init")
    //write http status
    w.WriteHeader(200)
}

//
// Handle Exit request
//
func ExitHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "{\"version\": \"mock-1.0.1\", \"exit\":\"OK\"}")
    log.Println("ExitHandler:Init")
    //write http status
    w.WriteHeader(200)
    os.Exit(3)
}
