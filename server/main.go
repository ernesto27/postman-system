package main

import (
	"fmt"
	"io"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Received headers:\n")

	// Iterate over headers and print them to the response and console
	for key, values := range r.Header {
		for _, value := range values {
			fmt.Fprintf(w, "%s: %s\n", key, value)
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	// Read and print body data
	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error reading body: %s\n", err)
		return
	}
	defer r.Body.Close()

	fmt.Fprintf(w, "\nBody:\n%s\n", string(bodyData))
	fmt.Printf("Body:\n%s\n", string(bodyData))
}

func main() {
	http.HandleFunc("/", handler)

	port := ":8080" // Define the port the server will run on
	fmt.Printf("Starting server on %s...\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
