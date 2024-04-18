package main

import (
	"log"
	"net/http"
)

func main1() {
	http.HandleFunc("/rolldice", rolldice1)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
