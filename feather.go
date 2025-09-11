package feather

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const VERSION string = "0.1.0"

type HandlerFunc func(c *Context)

type Route struct {
	Regex		*regexp.Regexp
	Params		[]string
	Handler		HandlerFunc
}

type Server struct {
	routes      map[string][]Route 				  // Stores the routes in a map [method]Route with the Handler, params name and the pattern stocked in Route
	middlewares []HandlerFunc                     // Stores the middlewares to run before the final routing method
}

func NewServer() *Server {
	return &Server{
		routes: make(map[string][]Route),
		middlewares: make([]HandlerFunc, 0),
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

	fragmentRegex := make([]string, 0)
	paramsList := make([]string, 0)

	// Get the different part of the path -> /:user/activate to [":user", "activate"]
	for fragment := range strings.SplitSeq(pattern, "/") { 
		parts := strings.Split(fragment, "|")
		
		if len(fragment) <= 0 {
			continue
		}

		if len(parts) == 1 && fragment[0] == ':' {
			// Default dynamic path /:user
			fragmentRegex = append(fragmentRegex, "([^/]+)")
			paramsList = append(paramsList, parts[0][1:])
		} else if len(parts) == 2 { 
			// Dynamic path with custom regex /:id|[0-9]+
			fragmentRegex = append(fragmentRegex, "(" + parts[1] + ")")
			paramsList = append(paramsList, parts[0][1:])
		} else { 
			// Static path
			fragmentRegex = append(fragmentRegex, regexp.QuoteMeta(fragment)) 
		} 
	}
	
	regexPattern := "^/" + strings.Join(fragmentRegex, "/") + "$"
	re, err := regexp.Compile(regexPattern)
	
	if err != nil {
		fmt.Printf("An error occured while parsing the dynamic route of \"%s\", the Regex isn't valid. \nFull error: %v\n", pattern, err)
		os.Exit(1)
	}

	route := Route{
		Regex: re,
		Params: paramsList,
		Handler: handler,
	}

	for _, method := range methods {
		if server.routes[method] == nil {
			server.routes[method] = make([]Route, 0)
		}

		server.routes[method] = append(server.routes[method], route)
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
	routes, ok := server.routes[reader.Method]
	if !ok {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	
	found := false
	index := -1
	params := make(map[string]string)
	
	for i, route := range routes {
		matches := route.Regex.FindStringSubmatch(reader.URL.Path)
		if len(matches) != 0 {
			found = true
			index = i
			
			for j, paramName := range route.Params {
				params[paramName] = matches[j + 1]
			}
			
			break
		} else {
			continue
		}
	}
	
	if !found {
		http.NotFound(writer, reader)
		return
	}

	context := &Context{
		Writer:  writer,
		Request: reader,
		Data:    make(map[string]any),
		Params:  params,
	}
	context.Data["PostFunc"] = make([]HandlerFunc, 0)
	context.Data["Abort"] = false

	for _, mw := range server.middlewares {
		mw(context)

		if context.Get("Abort").(bool) {
			break
		}
	}

	routes[index].Handler(context)

	postFuncs, ok := context.Data["PostFunc"].( []HandlerFunc )
	if !ok {
		return
	}

	for _, fn := range postFuncs {
		fn(context)
	}
}

func (server *Server) Listen(addr string) error {
	return http.ListenAndServe(addr, server)
}
