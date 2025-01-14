package main

// PostmanCollection represents the root structure of a Postman collection
type PostmanCollection struct {
	Info     Info       `json:"info"`
	Items    []Item     `json:"item"`
	Variable []Variable `json:"variable"`
}

// Info contains collection metadata
type Info struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Schema      string `json:"schema"`
}

// Item represents a request in the collection
type Item struct {
	Name     string     `json:"name"`
	Request  Request    `json:"request"`
	Response []Response `json:"response,omitempty"`
	Event    []Event    `json:"event,omitempty"`
	Items    []Item     `json:"item,omitempty"` // For folders/nested items
}

// Request contains the HTTP request details
type Request struct {
	Method string   `json:"method"`
	Header []Header `json:"header,omitempty"`
	URL    URL      `json:"url"`
	Body   Body     `json:"body,omitempty"`
}

// URL represents request URL structure
type URL struct {
	Raw      string   `json:"raw"`
	Protocol string   `json:"protocol,omitempty"`
	Host     []string `json:"host,omitempty"`
	Path     []string `json:"path,omitempty"`
}

// Header represents HTTP headers
type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type,omitempty"`
}

// Body represents request body
type Body struct {
	Mode string      `json:"mode,omitempty"`
	Raw  string      `json:"raw,omitempty"`
	Form interface{} `json:"formdata,omitempty"`
}

// Response represents a saved response
type Response struct {
	Name   string   `json:"name"`
	Status string   `json:"status"`
	Code   int      `json:"code"`
	Header []Header `json:"header"`
	Body   string   `json:"body"`
}

// Event represents pre-request scripts or tests
type Event struct {
	Listen string `json:"listen"`
	Script Script `json:"script"`
}

// Script contains script details
type Script struct {
	Type string   `json:"type"`
	Exec []string `json:"exec"`
}

// Variable represents collection variables
type Variable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}
