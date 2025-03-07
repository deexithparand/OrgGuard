package main

import (
	"fmt"
	"log"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("At good health"))
}

func main() {
	var (
		addr string = "localhost:"
		port string = "8080"
	)

	// routes
	http.HandleFunc("/", healthHandler)
	http.HandleFunc("/health", healthHandler)

	// started server
	fmt.Println("Server listening at -> ", addr+port)

	// currently mus is not handled
	err := http.ListenAndServe(addr+port, nil)
	if err != nil {
		log.Fatalf("error listening to server : %v", err)
		return
	}

	fmt.Println("Server says byee..")
}
