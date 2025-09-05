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
	http.ResponseWriter
	status int
}

func (recorder *responseRecorder) WriteHeader(code int) {
	recorder.status = code
	recorder.ResponseWriter.WriteHeader(code)
}

func Logging() feather.HandlerFunc {
	_, filepath, line, _ := runtime.Caller(1)
	file := strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))-1]
	fileName := strings.Split(file, ".")[0]
	
	date := time.Now()
	fmt.Printf("\033[1m%s\033[0m |\033[44m %s \033[0m| %-20s | %s\n",
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

func getStatusColor(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "\033[42m" // Vert
	case statusCode >= 300 && statusCode < 400:
		return "\033[45m" // Rose
	case statusCode >= 400 && statusCode < 500:
		return "\033[41m" // Rouge
	default:
		return "\033[43m" // Jaune
	}
}
func getMethodColor(method string) string {
	switch method {
	case "GET":
		return "\033[34m" // Bleu
	case "POST":
		return "\033[32m" // Cyan
	case "PUT":
		return "\033[33m" // Vert
	case "DELETE":
		return "\033[31m" // Rouge
	case "PATCH":
		return "\033[36m" // Magenta
	case "OPTIONS":
		return "\033[35m" // Jaune
	default:
		return "\033[37m" // Blanc
	}
}