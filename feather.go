package feather

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const VERSION string = "0.1.0"

// HandlerFunc represents a function that handles an HTTP request.
// It takes a single parameter:
//   - c: A pointer to the Context, which contains information about the HTTP request, response, and other data.
type HandlerFunc func(c *Context)

type Route struct {
	Regex *regexp.Regexp		// Regex is the compiled regular expression used to match the incoming request URL.
	Params []string 			// Params is a list of parameter names extracted from the dynamic segments of the route.
	Handler HandlerFunc 		// Handler is the function that will be executed when the route is matched.
}

type Server struct {
	// routes is a map where the key is the HTTP method (e.g., "GET", "POST") and the value is a slice of Route.
	// Each Route contains the compiled regular expression for matching the URL, the parameter names extracted from the route,
	// and the handler function to execute when the route is matched.
	Routes map[string][]Route

	// middlewares is a slice of HandlerFunc that represents middleware functions.
	// These functions are executed in the order they are added, before the final route handler is called.
	Middlewares []HandlerFunc
}

// NewServer creates and initializes a new instance of the Server struct.
//
// This function sets up the Server with an empty map for routes and an empty slice for middlewares.
// The Server is used to define routes, add middleware, and handle HTTP requests.
//
// Returns:
//   - *Server: A pointer to the newly created Server instance.
func NewServer() *Server {
	return &Server{
		Routes: make(map[string][]Route),
		Middlewares: make([]HandlerFunc, 0),
	}
}

// AddMiddleware appends one or more middleware functions to the server's middleware stack.
//
// Middleware functions are executed in the order they are added, before the final route handler is called.
//
// Parameters:
//   - middlewares: A variadic parameter of type HandlerFunc. Each HandlerFunc represents a middleware function
//     that takes a pointer to the Context as its argument. These middleware functions can modify the request,
//     response, or context data, and can also decide whether to abort the request processing.
//
// Returns:
//   - This function does not return any value.
func (server *Server) AddMiddleware(middlewares ...HandlerFunc) {
	for _, mw := range middlewares {
		server.Middlewares = append(server.Middlewares, mw)
	}
}

/*
	Handle registers a new route with the server, associating it with a specific URL pattern, handler function, 
	and one or more HTTP methods.

	The function supports dynamic URL segments, which can be defined using a colon (e.g., `/:user`). 
	Custom regular expressions can also be specified for dynamic segments (e.g., `/:id|[0-9]+`).

	Parameters:
			- pattern (string): The URL pattern for the route. It can include static segments, dynamic segments, 
					and optional custom regular expressions for dynamic segments.
			- handler (HandlerFunc): The function to execute when the route is matched. It receives a pointer 
					to the Context, which contains request and response data.
			- methods ([]string): A slice of HTTP methods (e.g., "GET", "POST") for which this route should be registered. 
					If no methods are provided, the default is ["GET"].

	Returns:
			- This function does not return any value.
*/
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
		if server.Routes[method] == nil {
			server.Routes[method] = make([]Route, 0)
		}

		server.Routes[method] = append(server.Routes[method], route)
	}
}

/*
	GET registers a new route with the HTTP method "GET" and associates it with a specific URL pattern and handler function.

	This function is a shorthand for calling the Handle method with the "GET" HTTP method. It allows you to define
	routes that respond to GET requests, which are typically used for retrieving data or resources.

	Parameters:
		- pattern (string): The URL pattern for the route. It can include static segments, dynamic segments, 
			and optional custom regular expressions for dynamic segments.
		- handler (HandlerFunc): The function to execute when the route is matched. It receives a pointer 
			to the Context, which contains request and response data.

	Returns:
		- This function does not return any value.
*/
func (server *Server) GET(pattern string, handler HandlerFunc) {
	server.Handle(pattern, handler, []string{"GET"})
}

/*
	POST registers a new route with the HTTP method "POST" and associates it with a specific URL pattern and handler function.

	This function is a shorthand for calling the Handle method with the "POST" HTTP method. It allows you to define
	routes that respond to POST requests, which are typically used for retrieving data or resources.

	Parameters:
		- pattern (string): The URL pattern for the route. It can include static segments, dynamic segments, 
			and optional custom regular expressions for dynamic segments.
		- handler (HandlerFunc): The function to execute when the route is matched. It receives a pointer 
			to the Context, which contains request and response data.

	Returns:
		- This function does not return any value.
*/
func (server *Server) POST(pattern string, handler HandlerFunc) {
	server.Handle(pattern, handler, []string{"POST"})
}

/*
	PUT registers a new route with the HTTP method "PUT" and associates it with a specific URL pattern and handler function.

	This function is a shorthand for calling the Handle method with the "PUT" HTTP method. It allows you to define
	routes that respond to PUT requests, which are typically used for retrieving data or resources.

	Parameters:
		- pattern (string): The URL pattern for the route. It can include static segments, dynamic segments, 
			and optional custom regular expressions for dynamic segments.
		- handler (HandlerFunc): The function to execute when the route is matched. It receives a pointer 
			to the Context, which contains request and response data.

	Returns:
		- This function does not return any value.
*/
func (server *Server) PUT(pattern string, handler HandlerFunc) {
	server.Handle(pattern, handler, []string{"PUT"})
}

/*
	PATCH registers a new route with the HTTP method "PATCH" and associates it with a specific URL pattern and handler function.

	This function is a shorthand for calling the Handle method with the "PATCH" HTTP method. It allows you to define
	routes that respond to PATCH requests, which are typically used for retrieving data or resources.

	Parameters:
		- pattern (string): The URL pattern for the route. It can include static segments, dynamic segments, 
			and optional custom regular expressions for dynamic segments.
		- handler (HandlerFunc): The function to execute when the route is matched. It receives a pointer 
			to the Context, which contains request and response data.

	Returns:
		- This function does not return any value.
*/
func (server *Server) PATCH(pattern string, handler HandlerFunc) {
	server.Handle(pattern, handler, []string{"PATCH"})
}

/*
	DELETE registers a new route with the HTTP method "DELETE" and associates it with a specific URL pattern and handler function.

	This function is a shorthand for calling the Handle method with the "DELETE" HTTP method. It allows you to define
	routes that respond to DELETE requests, which are typically used for retrieving data or resources.

	Parameters:
		- pattern (string): The URL pattern for the route. It can include static segments, dynamic segments, 
			and optional custom regular expressions for dynamic segments.
		- handler (HandlerFunc): The function to execute when the route is matched. It receives a pointer 
			to the Context, which contains request and response data.

	Returns:
		- This function does not return any value.
*/
func (server *Server) DELETE(pattern string, handler HandlerFunc) {
	server.Handle(pattern, handler, []string{"DELETE"})
}

/*
	ServeHTTP is the main entry point for handling HTTP requests in the Server.

	This function matches incoming HTTP requests against the registered routes based on the HTTP method and URL pattern.
	If a matching route is found, it creates a Context object, executes middleware functions, and invokes the route's handler.
	If no matching route is found, it responds with a 404 Not Found status. If the HTTP method is not allowed, it responds
	with a 405 Method Not Allowed status.

	Parameters:
		- writer (http.ResponseWriter): The HTTP response writer used to send data back to the client.
		- reader (*http.Request): The HTTP request object containing details about the incoming request.

	Returns:
		- This function does not return any value. It writes the HTTP response directly to the writer.
*/
func (server *Server) ServeHTTP(writer http.ResponseWriter, reader *http.Request) {
	routes, ok := server.Routes[reader.Method]
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

	for _, mw := range server.Middlewares {
		mw(context)

		if context.Get("Abort").(bool) {
			break
		}
	}

	routes[index].Handler(context)

	postFuncs, ok := context.Data["PostFunc"].([]HandlerFunc)
	if !ok {
		return
	}

	for _, fn := range postFuncs {
		fn(context)
	}
}

/*
	Listen starts the HTTP server on the specified address and begins handling incoming requests.

	This function uses the http.ListenAndServe function from the net/http package to bind the server
	to the given address and listen for incoming HTTP requests. The Server instance is used as the
	handler for these requests, routing them to the appropriate middleware and route handlers.

	Parameters:
		- addr (string): The address to listen on, in the format "host:port" (e.g., ":8080" for all
				interfaces on port 8080, or "127.0.0.1:8080" for localhost only).

	Returns:
		- error: If the server fails to start or encounters an error, this function returns the error.
				Otherwise, it blocks indefinitely and does not return.
*/
func (server *Server) Listen(addr string) error {
	return http.ListenAndServe(addr, server)
}
