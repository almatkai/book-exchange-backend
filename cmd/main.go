package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World 3!")
}

func main() {
    http.HandleFunc("/", handler)
    fmt.Println("Server running on port 8000")
    http.ListenAndServe(":8000", nil)
}
