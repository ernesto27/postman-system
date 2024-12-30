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

	// Print query parameters
	queryParams := r.URL.Query()
	if len(queryParams) > 0 {
		fmt.Fprintf(w, "\nQuery Parameters:\n")
		fmt.Printf("\nQuery Parameters:\n")
		for key, values := range queryParams {
			for _, value := range values {
				fmt.Fprintf(w, "%s: %s\n", key, value)
				fmt.Printf("%s: %s\n", key, value)
			}
		}
	}

	// Parse form data
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max memory
		fmt.Printf("Error parsing form: %s\n", err)
	}

	if r.MultipartForm != nil {
		fmt.Fprintf(w, "\nForm Data:\n")
		fmt.Printf("\nForm Data:\n")

		for key, values := range r.MultipartForm.Value {
			for _, value := range values {
				fmt.Fprintf(w, "Form field %s: %s\n", key, value)
				fmt.Printf("Form field %s: %s\n", key, value)
			}
		}

		for key, files := range r.MultipartForm.File {
			for _, file := range files {
				fmt.Fprintf(w, "File field %s: %s (%d bytes)\n", key, file.Filename, file.Size)
				fmt.Printf("File field %s: %s (%d bytes)\n", key, file.Filename, file.Size)
			}
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
