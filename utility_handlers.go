package main

import "net/http"

func ServerRunningHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Budget buddy server is running"))
}
