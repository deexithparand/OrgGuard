package handlers

import "net/http"

func GreetHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi Server"))
}
