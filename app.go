package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// FileData represents a file to be uploaded
type FileData struct {
	FieldName string `json:"fieldName"`
	FileName  string `json:"fileName"`
	Content   string `json:"content"`
}

// RequestParams represents the parameters for making an HTTP request
type RequestParams struct {
	Method   string            `json:"method"`
	URL      string            `json:"url"`
	Headers  map[string]string `json:"headers"`
	Body     string            `json:"body"`
	FormData map[string]string `json:"formData"`
	Files    []FileData        `json:"files"`
	BodyType string            `json:"bodyType"`
}

// App struct
type App struct {
	ctx    context.Context
	dynamo *DynamoDB
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	dynamo := NewDynamoDB("", "", "", "")
	a.ctx = ctx
	a.dynamo = dynamo
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetCollections() ([]Collection, error) {
	fmt.Println("Getting collections")
	collections, err := a.dynamo.GetCollections()
	if err != nil {
		return nil, err
	}

	return collections, nil
}

func (a *App) CreateCollection(name string, userID string) error {
	fmt.Println("Creating collection")
	err := a.dynamo.CreateCollection(name, userID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// MakeRequest handles HTTP requests from the frontend
func (a *App) MakeRequest(params RequestParams) (map[string]interface{}, error) {
	client := &http.Client{}
	startTime := time.Now()

	var reqBody io.Reader

	if params.BodyType == "form-data" {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		for key, value := range params.FormData {
			if err := writer.WriteField(key, value); err != nil {
				return nil, fmt.Errorf("failed to write form field: %v", err)
			}
		}

		for _, file := range params.Files {
			part, err := writer.CreateFormFile(file.FieldName, file.FileName)
			if err != nil {
				return nil, fmt.Errorf("failed to create form file: %v", err)
			}

			decoded, err := base64.StdEncoding.DecodeString(file.Content)
			if err != nil {
				return nil, fmt.Errorf("failed to decode file content: %v", err)
			}

			if _, err := part.Write(decoded); err != nil {
				return nil, fmt.Errorf("failed to write file content: %v", err)
			}
		}

		if err := writer.Close(); err != nil {
			return nil, fmt.Errorf("failed to close multipart writer: %v", err)
		}

		reqBody = body
		if params.Headers == nil {
			params.Headers = make(map[string]string)
		}
		params.Headers["Content-Type"] = writer.FormDataContentType()
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

	elapsed := time.Since(startTime)

	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Prepare response headers
	responseHeaders := make(map[string]string)
	headerSize := 0
	for key, values := range resp.Header {
		if len(values) > 0 {
			responseHeaders[key] = values[0]
			// Calculate header size: key length + ": " + value length + "\r\n"
			headerSize += len(key) + 2 + len(values[0]) + 2
		}
	}

	return map[string]interface{}{
		"body":          string(body),
		"headers":       responseHeaders,
		"status":        resp.Status,
		"responseSize":  len(body) + headerSize,
		"timeInSeconds": elapsed.Seconds(),
	}, nil
}
