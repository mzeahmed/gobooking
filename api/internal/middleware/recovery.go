package middleware

import (
	"log"
	"net/http"

	"github.com/mzeahmed/gobooking/internal/response"
)

// Recovery returns a middleware that recovers from panics raised by next,
// logs them and writes a generic 500 Internal Server Error JSON response
// instead of letting the panic crash the server.
func Recovery(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {

			if err := recover(); err != nil {

				log.Printf("panic: %v", err)

				response.JSON(w, http.StatusInternalServerError, map[string]string{
					"error": "internal server error",
				})
			}

		}()

		next.ServeHTTP(w, r)

	})
}
