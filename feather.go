package feather

import (
	"net/http"
)

type HandlerFunc func(c *Context)

type Server struct {
	routes      map[string]map[string]HandlerFunc // Stores the routes in a map [pattern][method]handler
	middlewares []HandlerFunc                     // Stores the middlewares to run before the final routing method
}

func NewServer() *Server {
	return &Server{
		routes: make(map[string]map[string]HandlerFunc),
	}
}

func (server *Server) AddMiddleware(middlewares ...HandlerFunc) {
	for _, mw := range middlewares {
		server.middlewares = append(server.middlewares, mw)
	}
}

func (server *Server) Handle(pattern string, handler HandlerFunc, methods []string) {
	if len(methods) == 0 {
		methods = []string{"GET"}
	}

	for _, method := range methods {
		if server.routes[pattern] == nil {
			server.routes[pattern] = make(map[string]HandlerFunc)
		}

		server.routes[pattern][method] = handler
	}
}

func (server *Server) GET(pattern string, handler HandlerFunc) {
	server.Handle(pattern, handler, []string{"GET"})
}
func (server *Server) POST(pattern string, handler HandlerFunc) {
	server.Handle(pattern, handler, []string{"POST"})
}
func (server *Server) PUT(pattern string, handler HandlerFunc) {
	server.Handle(pattern, handler, []string{"PUT"})
}
func (server *Server) PATCH(pattern string, handler HandlerFunc) {
	server.Handle(pattern, handler, []string{"PATCH"})
}
func (server *Server) DELETE(pattern string, handler HandlerFunc) {
	server.Handle(pattern, handler, []string{"DELETE"})
}

func (server *Server) ServeHTTP(writer http.ResponseWriter, reader *http.Request) {
	methods := server.routes[reader.URL.Path]
	if methods == nil {
		http.NotFound(writer, reader)
		return
	}

	handler, ok := server.routes[reader.URL.Path][reader.Method]
	if !ok {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	context := &Context{
		Writer:  writer,
		Request: reader,
		Data:    make(map[string]any),
	}

	for _, mw := range server.middlewares {
		mw(context)

		if context.Get("Abort").(bool) {
			break
		}
	}

	handler(context)
}

func (server *Server) Listen(addr string) error {
	return http.ListenAndServe(addr, server)
}
