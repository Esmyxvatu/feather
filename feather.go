package feather

import (
    "encoding/json"
    "net/http"
)

type Context struct {
    Writer  http.ResponseWriter
    Request *http.Request
    Params  map[string]string // pour les routes dynamiques
    Data    map[string]any    // stockage interne pour middleware
}

// Helper pour renvoyer du JSON
func (c *Context) JSON(status int, obj any) {
    c.Writer.Header().Set("Content-Type", "application/json")
    c.Writer.WriteHeader(status)
    json.NewEncoder(c.Writer).Encode(obj)
}

// Helper pour renvoyer du texte
func (c *Context) String(status int, s string) {
    c.Writer.Header().Set("Content-Type", "text/plain")
    c.Writer.WriteHeader(status)
    c.Writer.Write([]byte(s))
}

// Récupérer query param
func (c *Context) Query(key string) string {
    return c.Request.URL.Query().Get(key)
}

type HandlerFunc func(c *Context)

type Server struct {
    mux *http.ServeMux
}

func NewServer() *Server {
    return &Server{mux: http.NewServeMux()}
}

func (s *Server) Handle(pattern string, handler HandlerFunc) {
    s.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
        c := &Context{
            Writer:  w,
            Request: r,
            Data:    make(map[string]any),
        }
        handler(c)
    })
}

func (s *Server) Listen(addr string) error {
    return http.ListenAndServe(addr, s.mux)
}
