package middleware

import (
	"log"
	"net/http"

	"github.com/mzeahmed/go-booking/internal/response"
)

func Recovery(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {

			if err := recover(); err != nil {

				log.Printf("panic: %v", err)

				response.JSON(
					w,
					http.StatusInternalServerError,
					map[string]string{
						"error": "internal server error",
					},
				)
			}

		}()

		next.ServeHTTP(w, r)

	})
}
