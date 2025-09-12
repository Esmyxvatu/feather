# Feather

Feather is a lightweight, modular HTTP web framework for Go, designed for rapid development of RESTful APIs and web applications. It provides a simple routing system, middleware support, and convenient helpers for handling requests and responses.

## Features

- **Routing**: Register routes with static and dynamic segments, including custom regex for parameters.
- **Middleware**: Add global middleware functions for logging, CORS, authentication, etc.
- **Context**: Unified request/response context with helpers for JSON, HTML, files, headers, cookies, and more.
- **Extensible**: Easily add custom middlewares and handlers.
- **Minimal Dependencies**: Built on Go's standard library.

## Getting Started

### Installation

Add Feather to your Go project:

```bash
go get github.com/esmyxvatu/feather
```

### Basic Usage

```go
package main

import (
    "github.com/esmyxvatu/feather"
    "github.com/esmyxvatu/feather/middlewares"
)

func main() {
    server := feather.NewServer()

    // Add logging and CORS middleware
    server.AddMiddleware(
        middlewares.Logging(),
        middlewares.CORS(
            []string{"*"},
            []string{"GET", "POST", "PUT", "DELETE"},
            []string{"Content-Type"},
        ),
    )

    // Define routes
    server.GET("/", func(c *feather.Context) {
        c.String(200, "Welcome to Feather!")
    })

    server.GET("/user/:id|[0-9]+", func(c *feather.Context) {
        userID := c.Params["id"]
        c.JSON(200, map[string]string{"user_id": userID})
    })

    // Start server
    server.Listen(":8080")
}
```

## Routing

- Static routes: `/about`
- Dynamic routes: `/user/:id`
- Dynamic with regex: `/post/:slug|[a-z0-9\-]+`
- Wildcard: `/files/*path`

## Middleware

Middlewares are functions that run before the route handler. Use `AddMiddleware` to register them globally.

Example: Logging and CORS are included in `middlewares/`.

## Context Helpers

- `c.JSON(status, obj)` – Send JSON response
- `c.String(status, text)` – Send plain text
- `c.HTML(status, html)` – Send HTML
- `c.File(status, path)` – Send file
- `c.Status(status)` – Send status code only
- `c.Redirect(status, url)` – Redirect
- `c.SetHeader(key, value)` – Set response header
- `c.SetCookie(cookie)` – Set cookie
- `c.Query(key)` – Get query param
- `c.JSONBody(v)` – Parse JSON body
- `c.FormValue(key)` – Get form value

## License

GNU General Public License v3.0