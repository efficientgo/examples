package main

import (
	"net/http"
)

// Very minimal code for starting a web server. NOT production ready - always check errors and avoid globals (:
// Read more in "Efficient Go"; Example 2-6.

func handle(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("It kind of works!"))
}

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(handle))
}
