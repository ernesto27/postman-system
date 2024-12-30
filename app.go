package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// RequestParams represents the parameters for making an HTTP request
type RequestParams struct {
	Method   string            `json:"method"`
	URL      string            `json:"url"`
	Headers  map[string]string `json:"headers"`
	Body     string            `json:"body"`
	FormData map[string]string `json:"formData"`
	BodyType string            `json:"bodyType"`
}

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
func (a *App) MakeRequest(params RequestParams) (map[string]interface{}, error) {
	client := &http.Client{}

	var reqBody io.Reader

	if params.BodyType == "form-data" && len(params.FormData) > 0 {
		form := url.Values{}
		for key, value := range params.FormData {
			form.Add(key, value)
		}
		reqBody = strings.NewReader(form.Encode())
	} else if params.BodyType == "raw" {
		reqBody = strings.NewReader(params.Body)
	}

	req, err := http.NewRequest(params.Method, params.URL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	for key, value := range params.Headers {
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
