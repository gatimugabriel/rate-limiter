package middleware

import (
	"fmt"
	"github.com/gabrielgatimu/rate-limiter/bucket"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	clients map[string]*bucket.Client
	Mu      sync.Mutex
)

func init() {
	clients = make(map[string]*bucket.Client)
}

func RateLimit(next http.Handler) http.Handler {
	//---@routine -> kill old clients
	go func() {
		for {
			time.Sleep(1 * time.Second)
			Mu.Lock()

			for ip, client := range clients {
				if time.Since(client.LastSeen) > 1*time.Minute {
					fmt.Println("\t Killed client...", ip)
					delete(clients, ip)
				}
			}
			Mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, ok := r.Context().Value("ip").(string)
		if !ok {
			log.Fatal("Missing IP")
		}

		Mu.Lock()
		defer Mu.Unlock()

		var client *bucket.Client
		if value, exists := clients[ip]; !exists {
			clients[ip] = &bucket.Client{
				Bucket: *bucket.CreateBucket(1, 1),
			}
			client = clients[ip]

			// ---@routine -> start the bucket refill ()
			//go client.Bucket.StartRefillBucket()
		} else {
			client = value
		}
		client.LastSeen = time.Now()

		if client.Bucket.Tokens == 0 {
			fmt.Println(w, "Too many requests for IP: ", ip, "try after a minute")

			w.WriteHeader(http.StatusTooManyRequests)
			_, err := fmt.Fprintf(w, "Too many requests for IP: %s", ip)
			if err != nil {
				return
			}
			return
		} else {
			// Decrement token for allowed request
			client.Bucket.Tokens--
		}

		next.ServeHTTP(w, r)
	})
}
