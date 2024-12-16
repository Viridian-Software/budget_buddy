package main

import "net/http"

func (cfg *apiConfig) ServerRunningHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Budget buddy server is running"))
}
