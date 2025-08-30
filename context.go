package feather

import (
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

type Context struct {
    Writer  http.ResponseWriter	// Default net/http writer object
    Request *http.Request		// Default net/http request object
    Params  map[string]string 	// Map storing the data of dynamic routing
    Data    map[string]any    	// Internal stockage for middleware
}

//==================================================== Helper for the response ==========================================================================================

// 
func (c *Context) JSON(status int, obj any) {
    c.Writer.Header().Set("Content-Type", "application/json")
    c.Writer.WriteHeader(status)

    json.NewEncoder(c.Writer).Encode(obj)
}

// Helper pour envoyer du texte pur
func (c *Context) String(status int, s string) {
    c.Writer.Header().Set("Content-Type", "text/plain")
    c.Writer.WriteHeader(status)

    c.Writer.Write([]byte(s))
}

// Helper pour envoyer du contenu HTML
func (c *Context) HTML(status int, content string) {
    c.Writer.Header().Set("Content-Type", "text/html")
    c.Writer.WriteHeader(status)

    c.Writer.Write([]byte(content))
}

// Helper pour envoyer un fichier
func (c *Context) File(status int, path string) {
	file, err := os.Open(path)
	if err != nil { http.Error(c.Writer, "File not found", http.StatusNotFound) }
	
	defer file.Close()
	
	extension := filepath.Ext(path)
	ctype := mime.TypeByExtension(extension)
	if ctype == "" { ctype = "application/octet-stream" } // Fallback
	
	c.Writer.Header().Set("Content-Type", ctype)
	c.Writer.WriteHeader(status)
	
	io.Copy(c.Writer, file)
}

// Helper pour renvoyer uniquement un status. Utile pour ping
func (c *Context) Status(status int) {
	c.Writer.WriteHeader(status)
	c.Writer.Write([]byte{})
}

// Helper pour rediriger le client sur une autre URL
func (c *Context) Redirect(status int, url string) {
	http.Redirect(c.Writer, c.Request, url, status)
}

// Helper pour renvoyer une erreur
func (c *Context) Error(status int, message string) {
	http.Error(c.Writer, message, status)
}

// Helper pour modifier un header manuellement
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Add(key, value)
}

// Helper pour définir le type du contenu manuellement
func (c *Context) ContentType(value string) {
	c.Writer.Header().Set("Content-Type", value)
}

// Helper pour ajouter un cookie à la réponse
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Writer, cookie)
}

//==================================================== Helper for the request ===========================================================================================

// Helper pour récupérer un Query param. Correspond au paramètres fourni dans la requête ( /foo?bar=pee )
func (c *Context) Query(key string) string {
    return c.Request.URL.Query().Get(key)
}

/*
Parse the body data into JSON format. This method reads the request body and decodes it into a structure using encoding/json.
Return an error if something wrong happened and fill the given structure with the data of the body request.
*/
func (c *Context) JSONBody(v any) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil { return err }
	
	err = json.Unmarshal(body, v)
	if err != nil { return err }

	return nil
}

// ...
func (c *Context) Header(key string) string {
	return c.Request.Header.Get(key)
}

// Helper pour récupérer la valeur d'un cookie
func (c *Context) Cookie(name string) (*http.Cookie, error) {
	return c.Request.Cookie(name)
}

/* 
Parse the request's form data. This method reads the request body and decodes it into key-value pairs, according to the Content-Type header (usually application/x-www-form-urlencoded).
Return the value associated with the given key from the form data. If the key is not present, return an empty string.
*/
func (c *Context) FormValue(key string) string {
	err := c.Request.ParseForm()
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
	}

	return c.Request.FormValue(key)
}

//==================================================== Helper for middlewares ===========================================================================================

func (c *Context) Set(key string, value any) {
	c.Data[key] = value
}

func (c *Context) Get(key string) any {
	return c.Data[key]
}

func (c *Context) ClientIP() string {
	return c.Request.RemoteAddr
}