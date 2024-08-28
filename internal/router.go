package internal

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

const (
	TemplatePath = "internal/templates/response.txt"
)

// ResponseData represents the structure of the data sent to the template
type ResponseData struct {
	Route    string
	Method   string
	Headers  map[string]string
	OriginIP string
}

// StartServer initializes and starts the HTTP server
func StartServer() {
	// Load the template
	tmpl, err := LoadTemplate(TemplatePath)
	if err != nil {
		log.Fatalf("Failed to load template: %v", err)
	}

	// Define the route and handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HandleRequest(w, r, tmpl)
	})

	// Start the server
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// LoadTemplate loads the text template from the specified file path
func LoadTemplate(path string) (*template.Template, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Combine the current working directory with the relative path
	absolutePath := filepath.Join(cwd, path)

	tmpl, err := template.ParseFiles(filepath.Clean(absolutePath))
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

// HandleRequest processes the incoming request and responds using the template
func HandleRequest(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {

	// Set Content-Type to text/plain
	w.Header().Set("Content-Type", "text/plain")

	// Collect headers and other request information
	data := CollectRequestData(r)

	responseCode := 200
	// Set the custom response status code if present in header
	if code, ok := data.Headers["X-Response-Status"]; ok {
		if statusCode, err := strconv.Atoi(code); err == nil {
			responseCode = statusCode
		}
	}

	w.WriteHeader(responseCode)

	// Log received request
	log.Default().Printf("%v %v: %v", responseCode, r.Method, r.URL.Path)

	// Render the template
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// CollectRequestData gathers all necessary information from the request
func CollectRequestData(r *http.Request) ResponseData {
	headers := make(map[string]string)
	for name, values := range r.Header {
		headers[name] = strings.Join(values, ", ")
	}

	return ResponseData{
		Route:    r.URL.Path,
		Method:   r.Method,
		Headers:  headers,
		OriginIP: r.RemoteAddr,
	}
}
