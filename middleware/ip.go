package middleware

import (
	"context"
	rip "github.com/vikram1565/request-ip"
	"log"
	"net/http"
)

func IPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		incomingIP := rip.GetClientIP(r)
		log.Println("Client IP:", incomingIP)

		// store the IP in request context
		newCtx := context.WithValue(r.Context(), "ip", incomingIP)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}
