package main

import (
	"log"
	"net/http"
)

const url = "https://localhost:8080"

func main() {
	_, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
}
