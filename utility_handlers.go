package main

import "net/http"

func (cfg *apiConfig) ServerRunningHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Budget buddy server is running"))
}
