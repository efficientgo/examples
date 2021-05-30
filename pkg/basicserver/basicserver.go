package main

import (
	"net/http"
)

func handle(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("It kind of works!"))
}

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(handle))
}
