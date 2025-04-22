package handlers

import "net/http"

func TodosHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from handler"))
}
