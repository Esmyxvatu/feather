package middlewares

import (
	"net/http"
	"strings"

	"github.com/esmyxvatu/feather"
)

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
