package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// MakeRequest handles HTTP requests from the frontend
func (a *App) MakeRequest(
	method string,
	urlRequest string,
	headersRequest map[string]string,
	bodyRequest string,
	formData map[string]string,
	bodyType string,
) (map[string]interface{}, error) {
	client := &http.Client{}

	var reqBody io.Reader

	if bodyType == "form-data" && len(formData) > 0 {

		form := url.Values{}
		for key, value := range formData {
			form.Add(key, value)
		}
		reqBody = strings.NewReader(form.Encode())

	} else if bodyType == "raw" {
		reqBody = strings.NewReader(bodyRequest)
	}

	req, err := http.NewRequest(method, urlRequest, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	for key, value := range headersRequest {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Prepare response headers
	responseHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			responseHeaders[key] = values[0]
		}
	}

	return map[string]interface{}{
		"body":    string(body),
		"headers": responseHeaders,
		"status":  resp.Status,
	}, nil
}
