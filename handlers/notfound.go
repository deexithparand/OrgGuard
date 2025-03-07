package handlers

import "net/http"

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	http.Error(w, "\""+path+"\" : "+" Route Not Found", 404)
}
