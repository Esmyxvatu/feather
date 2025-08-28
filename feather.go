package feather

import (
    "net/http"
)

type HandlerFunc func(c *Context)

type Server struct {
    mux *http.ServeMux
}

func NewServer() *Server {
    return &Server{mux: http.NewServeMux()}
}

func (server *Server) Handle(pattern string, handler HandlerFunc) {
    server.mux.HandleFunc(pattern, func(writer http.ResponseWriter, reader *http.Request) {
        context := &Context{
            Writer:  writer,
            Request: reader,
            Data:    make(map[string]any),
        }
        handler(context)
    })
}

func (server *Server) Listen(addr string) error {
    return http.ListenAndServe(addr, server.mux)
}
