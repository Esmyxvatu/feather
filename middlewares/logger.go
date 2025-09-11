package middlewares

import (
	"net/http"
	"time"
	"fmt"
	"strings"
	"runtime"
	
	"github.com/esmyxvatu/feather"
)

type responseRecorder struct {
	/*
		ResponseWriter is an embedded field that allows the responseRecorder
		to act as an http.ResponseWriter. It is used to write the HTTP response
		to the client.
	*/
	http.ResponseWriter

	/*
		status is an integer field that records the HTTP status code
		of the response. It is used for logging and monitoring purposes.
	*/
	status int
}

/*
	WriteHeader sets the HTTP status code for the response and records it.

	Parameters:
	- code (int): The HTTP status code to be set for the response.

	Returns:
	- None
*/
func (recorder *responseRecorder) WriteHeader(code int) {
	recorder.status = code
	recorder.ResponseWriter.WriteHeader(code)
}

/*
	Logging is a middleware function that logs HTTP requests and responses in a structured format.
	It provides details such as the timestamp, HTTP status code, client IP, HTTP method, request path, and response time.

	Parameters:
	- None

	Returns:
	- feather.HandlerFunc: A function that can be used as middleware in a Feather application.
*/
func Logging() feather.HandlerFunc {
	_, filepath, line, _ := runtime.Caller(1)
	file := strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))-1]
	fileName := strings.Split(file, ".")[0]

	date := time.Now()
	fmt.Printf("\033[1m%s\033[0m │\033[44m %s \033[0m│ %-20s │ %s\n",
		date.Format("2006/01/02 15:04:05.000"),
		"DEBUG",
		fileName + ":" + fmt.Sprint(line),
		"Logger initialized, using Feather v" + feather.VERSION,
	)

	return func(c *feather.Context) {
		start := time.Now()
		status := http.StatusOK

		recorder := &responseRecorder{
			ResponseWriter: c.Writer,
			status: status,
		}

		c.Writer = recorder

		c.Post(
			func(*feather.Context) {
				duration := time.Since(start)
				duration = duration.Round(time.Millisecond)
				if duration < 0 {
					duration = 0
				}

				padding := (7 - len(fmt.Sprint(recorder.status))) / 2
				status := fmt.Sprintf("%s%s%s",
					strings.Repeat(" ", padding),
					fmt.Sprint(recorder.status),
					strings.Repeat(" ", 7-len(fmt.Sprint(recorder.status))-padding),
				)
				status = fmt.Sprintf("%s%s%s", getStatusColor(recorder.status), status, "\033[0m") // Color of the HTTP status
				method := fmt.Sprintf("%s%s%s", getMethodColor(c.Request.Method), c.Request.Method, "\033[0m")   // Color of the method

				// Show the log in the format wanted
				fmt.Printf("\033[1m%s\033[0m │%s│ %-20s │ %s '%s' \033[2m%s\033[0m\n",
					start.Format("2006/01/02 15:04:05.000"), // Date/Hour
					status,                                  // Code HTTP
					c.ClientIP(),                            // IP
					method,                                  // Method
					c.Request.URL.Path,                      // Path
					duration,                                // Duration
				)
			},
		)
	}
}

/*
	getStatusColor determines the appropriate ANSI color code for a given HTTP status code.

	Parameters:
	- statusCode (int): The HTTP status code for which the color is to be determined.

	Returns:
	- string: An ANSI color code representing the category of the HTTP status code.
*/
func getStatusColor(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "\033[42m" // Green
	case statusCode >= 300 && statusCode < 400:
		return "\033[45m" // Magenta
	case statusCode >= 400 && statusCode < 500:
		return "\033[41m" // Red
	default:
		return "\033[43m" // Yellow
	}
}

/*
	getMethodColor determines the appropriate ANSI color code for a given HTTP method.

	Parameters:
	- method (string): The HTTP method for which the color is to be determined.

	Returns:
	- string: An ANSI color code representing the HTTP method.
*/
func getMethodColor(method string) string {
	switch method {
	case "GET":
		return "\033[34m" // Blue
	case "POST":
		return "\033[32m" // Cyan
	case "PUT":
		return "\033[33m" // Green
	case "DELETE":
		return "\033[31m" // Red
	case "PATCH":
		return "\033[36m" // Magenta
	case "OPTIONS":
		return "\033[35m" // Yellow
	default:
		return "\033[37m" // White
	}
}