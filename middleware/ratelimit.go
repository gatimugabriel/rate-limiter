package middleware

import (
	"fmt"
	"log"
	"net/http"
)

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, ok := r.Context().Value("ip").(string)
		if !ok {
			log.Fatal("Missing IP")
		}

		fmt.Println("hi there", ip)

		next.ServeHTTP(w, r)
	})
}
