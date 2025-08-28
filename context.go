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
    Writer  http.ResponseWriter
    Request *http.Request
    Params  map[string]string // pour les routes dynamiques
    Data    map[string]any    // stockage interne pour middleware
}

//==================================================== Helper pour la réponse ===========================================================================================

// Helper pour envoyer du contenu sous format JSON.
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

//==================================================== Helper pour la requête ===========================================================================================

// Helper pour récupérer un Query param. Correspond au paramètres fourni dans la requête ( /foo?bar=pee )
func (c *Context) Query(key string) string {
    return c.Request.URL.Query().Get(key)
}

// Helper pour récupérer le contenu de la requête sous format JSON
func (c *Context) JSONBody(v any) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil { return err }
	
	err = json.Unmarshal(body, v)
	if err != nil { return err }

	return nil
}

// Helper pour récupérer la valeur d'un Header de la requête
func (c *Context) Header(key string) string {
	return c.Request.Header.Get(key)
}

// Helper pour récupérer la valeur d'un cookie
func (c *Context) Cookie(name string) (*http.Cookie, error) {
	return c.Request.Cookie(name)
}

//==================================================== Helper pour les middlewares ======================================================================================

func (c *Context) Set(key string, value any) {
	c.Data[key] = value
}