package main

import (
	"log"
	"net/http"
)

func main() {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(handle),
	}

	log.Printf("Serving on https://0.0.0.0:8080")
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}

func handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got connection: %s", r.Proto)

	if r.URL.Path == "/2nd" {
		log.Println("Handling 2nd")
		w.Write([]byte("Hello Again!"))
		return
	}

	log.Println("Handling 1st")
	pusher, ok := w.(http.Pusher)
	if !ok {
		log.Println("Can't push to client")
	} else {
		err := pusher.Push("/2nd", nil)
		if err != nil {
			log.Printf("Failed push: %v", err)
		}
	}

	w.Write([]byte("Hello"))
}
