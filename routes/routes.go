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

	r.HandleFunc("/unlimited", unlimitedHandler)

	r.HandleFunc("/limited", limitedHandler)

	r.Use(mid.IPMiddleware)
	r.Use(mid.RateLimit)
	r.Use(mid.RequestLogger)

	log.Println("server started on: http://127.0.1.1:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func unlimitedHandler(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintf(writer, "Unlimited! Let's Go!\n")
	if err != nil {
		return
	}
}

func limitedHandler(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintf(writer, "Limited, don't overuse me!\n")
	if err != nil {
		return
	}
}
