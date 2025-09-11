package middlewares

import (
	"net/http"
	"strings"

	"github.com/esmyxvatu/feather"
)

/*
CORS is a middleware function that sets Cross-Origin Resource Sharing (CORS) headers
on HTTP responses. It allows the server to specify which origins, methods, and headers
are permitted for cross-origin requests.

Parameters:
		- allowedOrigins: A slice of strings specifying the allowed origins.
		- allowedMethods: A slice of strings specifying the allowed HTTP methods.
		- allowedHeaders: A slice of strings specifying the allowed HTTP headers.

Returns:
		- A feather.HandlerFunc that applies the CORS headers to the HTTP response.
*/
func CORS(allowedOrigins []string, allowedMethods []string, allowedHeaders []string) feather.HandlerFunc {
	return func(c *feather.Context) {
		c.SetHeader("Access-Control-Allow-Origin", strings.Join(allowedOrigins, ","))
		c.SetHeader("Access-Control-Allow-Methods", strings.Join(allowedMethods, ","))
		c.SetHeader("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ","))

		if c.Request.Method == http.MethodOptions {
			c.Status(http.StatusOK)
			return
		}
	}
}
