package routes

import (
	"fmt"
	mid "github.com/gabrielgatimu/rate-limiter/middleware"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func HandleRequests() {
	r := mux.NewRouter()

	r.HandleFunc("/unlimited", func(writer http.ResponseWriter, request *http.Request) {
		_, err := fmt.Fprintf(writer, "Unlimited! Let's Go!")
		if err != nil {
			return
		}
	})

	r.HandleFunc("/limited", func(writer http.ResponseWriter, request *http.Request) {
		_, err := fmt.Fprintf(writer, "Limited, don't overuse me!")
		if err != nil {
			return
		}
	})

	r.Use(mid.IPMiddleware)
	r.Use(mid.RateLimit)

	log.Println("server started on: http://127.0.1.1:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
