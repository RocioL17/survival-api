package handlers

import (
 "fmt"
 "net/http"
)

// type Data struct {
//  ID   string `json:"id"`
//  Name string `json:"name"`
// }

func MakeCase(w http.ResponseWriter, r *http.Request) {
 w.Header().Set("Content-Type", "application/json")
 
}

func VerifyChoice(w http.ResponseWriter, r *http.Request) {
 w.Header().Set("Content-Type", "application/json")
 
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
 w.Header().Set("Content-Type", "text/plain")
 w.WriteHeader(http.StatusNotFound)
 fmt.Fprintln(w, "404 - Page not found")
}