package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from backend server on port: 8082")
	})

	fmt.Println("Backend server started on 8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
