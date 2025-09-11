package feather

import (
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"html/template"
	"math/rand"
)

// Context represents the state and data associated with an HTTP request and response.
// It provides methods for handling requests, sending responses, and storing data
// for middleware and handlers.
type Context struct {
    Writer  http.ResponseWriter // Writer is the HTTP response writer used to construct the HTTP response.
    Request *http.Request       // Request is the HTTP request object containing details about the client's request.
    Params  map[string]string   // Params is a map that stores dynamic route parameters extracted from the URL.
    Data    map[string]any      // Data is a map for storing arbitrary key-value pairs, typically used by middleware.
}

//==================================================== Helper for the response ==========================================================================================

// JSON sends a JSON-encoded response with the specified HTTP status code.
//
// Parameters:
//   - status: The HTTP status code to set for the response.
//   - obj: The object to be JSON-encoded and sent in the response body.
//
// This function sets the "Content-Type" header to "application/json",
// writes the HTTP status code to the response, and encodes the provided
// object as JSON into the response body.
func (c *Context) JSON(status int, obj any) {
    c.Writer.Header().Set("Content-Type", "application/json")
    c.Writer.WriteHeader(status)

    json.NewEncoder(c.Writer).Encode(obj)
}

// String sends a plain text response with the specified HTTP status code.
//
// Parameters:
//   - status: The HTTP status code to set for the response.
//   - s: The string content to be sent in the response body.
//
// This function sets the "Content-Type" header to "text/plain",
// writes the HTTP status code to the response, and writes the provided
// string content into the response body.
func (c *Context) String(status int, s string) {
    c.Writer.Header().Set("Content-Type", "text/plain")
    c.Writer.WriteHeader(status)

    c.Writer.Write([]byte(s))
}

// HTML sends an HTML response with the specified HTTP status code.
//
// Parameters:
//   - status: The HTTP status code to set for the response.
//   - content: The HTML content to be sent in the response body.
//
// This function sets the "Content-Type" header to "text/html",
// writes the HTTP status code to the response, and writes the provided
// HTML content into the response body.
func (c *Context) HTML(status int, content string) {
    c.Writer.Header().Set("Content-Type", "text/html")
    c.Writer.WriteHeader(status)

    c.Writer.Write([]byte(content))
}

// File sends the contents of a file as the HTTP response.
//
// Parameters:
//   - status: The HTTP status code to set for the response.
//   - path: The file system path of the file to be sent.
//
// This function determines the file's MIME type based on its extension,
// sets the "Content-Type" header accordingly, and writes the file's
// contents to the response body. If the file cannot be opened, it sends
// a "404 Not Found" error response.
func (c *Context) File(status int, path string) {
	file, err := os.Open(path)
	if err != nil {
		http.Error(c.Writer, "File not found", http.StatusNotFound)
	}

	defer file.Close()

	extension := filepath.Ext(path)
	ctype := mime.TypeByExtension(extension)
	if ctype == "" {
		ctype = "application/octet-stream" // Fallback
	}

	c.Writer.Header().Set("Content-Type", ctype)
	c.Writer.WriteHeader(status)

	io.Copy(c.Writer, file)
}

// Status sends an HTTP response with the specified status code and an empty body.
//
// Parameters:
//   - status: The HTTP status code to set for the response.
//
// This function sets the HTTP status code for the response and writes an empty body.
func (c *Context) Status(status int) {
	c.Writer.WriteHeader(status)
	c.Writer.Write([]byte{})
}

// Redirect sends an HTTP redirect response to the client.
//
// Parameters:
//   - status: The HTTP status code to set for the redirect response.
//             Common values include 301 (Moved Permanently) and 302 (Found).
//   - url: The target URL to which the client should be redirected.
//
// This function uses the http.Redirect method to send a redirect response
// with the specified status code and target URL.
func (c *Context) Redirect(status int, url string) {
	http.Redirect(c.Writer, c.Request, url, status)
}

// Error sends an HTTP error response with the specified status code and message.
//
// Parameters:
//   - status: The HTTP status code to set for the error response.
//   - message: The error message to be sent in the response body.
//
// This function uses the http.Error method to send an error response
// with the provided status code and message.
func (c *Context) Error(status int, message string) {
	http.Error(c.Writer, message, status)
}

// SetHeader adds a header to the HTTP response.
//
// Parameters:
//   - key: The name of the header to set.
//   - value: The value to associate with the header.
//
// This function adds the specified header and its value to the HTTP response.
// If the header already exists, the new value is appended to the existing values.
func (c *Context) SetHeader(key string, value string) {
	
	c.Writer.Header().Add(key, value)
}

// ContentType sets the "Content-Type" header for the HTTP response.
//
// Parameters:
//   - value: The MIME type to set as the "Content-Type" header value.
//
// This function updates the "Content-Type" header in the HTTP response
// to the specified MIME type, replacing any existing value.
func (c *Context) ContentType(value string) {
	c.Writer.Header().Set("Content-Type", value)
}

// SetCookie adds a Set-Cookie header to the HTTP response.
//
// Parameters:
//   - cookie: A pointer to an http.Cookie object that contains the
//             cookie's name, value, and other attributes such as
//             expiration, path, domain, etc.
//
// This function uses the http.SetCookie method to add the specified
// cookie to the HTTP response. It does not return any value.
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Writer, cookie)
}

// Template renders HTML templates with optional data and custom functions.
//
// Parameters:
//   - files: A slice of strings representing the file paths of the templates to be parsed.
//   - data: The data to be passed to the template for rendering. This can be any type.
//   - funcs: A template.FuncMap containing custom functions to be used within the templates.
//
// This function generates a random name for the template, parses the provided files,
// and executes the template with the given data. If an error occurs during execution,
// it sends an HTTP 500 Internal Server Error response with the error message.
func (c *Context) Template(files []string, data any, funcs template.FuncMap) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := "abcdefghijklmnopqrstuvwxyz"
	word := make([]byte, 32)

	for i := range 32 {
		word[i] = chars[r.Intn(len(chars))]
	}

	tmpl := template.Must(
		template.New(string(word)).Funcs(funcs).ParseFiles(files...),
	)

	err := tmpl.Execute(c.Writer, data)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
	}
}

//==================================================== Helper for the request ===========================================================================================

// Query retrieves the value of a query parameter from the URL.
//
// Parameters:
//   - key: The name of the query parameter to retrieve.
//
// Returns:
//   - The value of the specified query parameter as a string.
//     If the parameter is not present, it returns an empty string.
func (c *Context) Query(key string) string {
    return c.Request.URL.Query().Get(key)
}

// JSONBody reads the request body and unmarshals it into the provided structure.
//
// Parameters:
//   - v: A pointer to the structure where the JSON data will be unmarshaled.
//
// Returns:
//   - An error if reading the body or unmarshaling the JSON fails.
//     If successful, the provided structure is populated with the request data.
func (c *Context) JSONBody(v any) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil { return err }

	err = json.Unmarshal(body, v)
	if err != nil { return err }

	return nil
}

// Header retrieves the value of a specific request header.
//
// Parameters:
//   - key: The name of the header to retrieve.
//
// Returns:
//   - The value of the specified header as a string.
//     If the header is not present, it returns an empty string.
func (c *Context) Header(key string) string {
	return c.Request.Header.Get(key)
}

// Cookie retrieves a specific cookie from the HTTP request.
//
// Parameters:
//   - name: The name of the cookie to retrieve.
//
// Returns:
//   - A pointer to an http.Cookie object representing the cookie with the specified name.
//   - An error if the cookie is not found or if there is an issue retrieving it.
func (c *Context) Cookie(name string) (*http.Cookie, error) {
	return c.Request.Cookie(name)
}

// FormValue parses the request's form data and retrieves the value for the specified key.
//
// Parameters:
//   - key: The name of the form field to retrieve the value for.
//
// Returns:
//   - The value of the specified form field as a string.
//     If the form field is not present, it returns an empty string.
//     If there is an error parsing the form data, it sends an HTTP 400 Bad Request response
//     and does not return a value.
func (c *Context) FormValue(key string) string {
	err := c.Request.ParseForm()
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
	}

	return c.Request.FormValue(key)
}

//==================================================== Helper for middlewares ===========================================================================================

// Set stores a key-value pair in the Context's Data map. This method should only be used by middlewares.
//
// Parameters:
//   - key: A string representing the key under which the value will be stored.
//   - value: The value to be stored, which can be of any type.
//
// This function does not return any value. It updates the Context's Data map
// by associating the specified key with the provided value.
func (c *Context) Set(key string, value any) {
	c.Data[key] = value
}

// Get retrieves the value associated with the specified key from the Context's Data map. This method should only be used by middlewares.
//
// Parameters:
//   - key: A string representing the key whose associated value is to be retrieved.
//
// Returns:
//   - The value associated with the specified key, which can be of any type.
//     If the key does not exist in the Data map, it returns nil.
func (c *Context) Get(key string) any {
	return c.Data[key]
}

// ClientIP retrieves the IP address of the client making the request.
//
// This function does not take any parameters.
//
// Returns:
//   - A string representing the client's IP address as obtained from the
//     RemoteAddr field of the HTTP request.
func (c *Context) ClientIP() string {
	return c.Request.RemoteAddr
}

// Abort halts the execution of any subsequent middleware or handlers. This method should only be used by middlewares.
//
// This function sets the "Abort" key in the Context's Data map to true,
// signaling that the request processing should be stopped immediately.
// It does not take any parameters and does not return any value.
func (c *Context) Abort() {
	c.Data["Abort"] = true
}

// Post appends a new handler function to the "PostFunc" middleware chain stored in the Context's Data map. This method should only be used by middlewares.
//
// Parameters:
//   - function: A HandlerFunc representing the middleware or handler function to be added to the "PostFunc" chain.
//
// This function retrieves the existing "PostFunc" middleware chain from the Context's Data map,
// appends the provided handler function to the chain, and updates the "PostFunc" entry in the Data map.
// It does not return any value.
func (c *Context) Post(function HandlerFunc) {
	postMw := c.Data["PostFunc"]

	c.Data["PostFunc"] = append(postMw.([]HandlerFunc), function)
}
