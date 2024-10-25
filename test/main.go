package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	args := os.Args

	if len(args) <= 1 {
		log.Fatal("Port Number required")
	}

	port := args[1]

	originServerHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Printf("[origin server] received request at: %s\n", time.Now())
		fmt.Printf("[origin server] headers: %v\n\n", req.Header)
		rw.Header().Set("Test-Header", "Test-Header info")
		fmt.Fprintf(rw, "origin test server response from test server on %v for %v %v\nHeader:%v\n",
			port, req.Method, req.URL, req.Header)
	})

	fmt.Printf("Run Test Server on port %v\n", port)
	log.Fatal(http.ListenAndServe(":"+port, originServerHandler))
}
