package main

import (
	"github.com/kernle32dll/emissione-go"

	"log"
	"net/http"
)

// User is a just sample struct for showcasing.
type User struct {
	Name string `json:"name",xml:"Name"`
}

func main() {
	router := http.NewServeMux()

	em := emissione.Default()

	router.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		em.Write(w, r, http.StatusOK, User{Name: "Björn Gerdau"})
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}
