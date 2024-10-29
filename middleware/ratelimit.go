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
	once    sync.Once
)

func init() {
	clients = make(map[string]*bucket.Client)
}

func startClientCleanup() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		Mu.Lock()

		if len(clients) == 0 {
			Mu.Unlock()
			continue
		}

		//fmt.Println("Killer in action")
		for ip, client := range clients {
			if time.Since(client.LastSeen) > 5*time.Minute {
				fmt.Println("\t Killed client...", ip)
				delete(clients, ip)
			}
		}

		Mu.Unlock()
	}
}

func RateLimit(next http.Handler) http.Handler {
	// start cleanup(kills old clients) routine only once
	once.Do(func() { go startClientCleanup() })

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, ok := r.Context().Value("ip").(string)
		if !ok {
			log.Fatal("Missing IP")
		}

		Mu.Lock()
		defer Mu.Unlock()

		var client *bucket.Client
		fmt.Println("\tconnected clients: ", &clients)
		if value, exists := clients[ip]; !exists {
			clients[ip] = &bucket.Client{
				Bucket: *bucket.CreateBucket(10, 1),
			}
			client = clients[ip]

			// ---@routine -> start the bucket refill ()
			go client.Bucket.StartRefillBucket()
		} else {
			client = value
		}
		client.LastSeen = time.Now()

		if client.Bucket.Tokens == 0 {
			w.WriteHeader(http.StatusTooManyRequests)
			_, err := fmt.Fprintf(w, "Too many requests/second for IP: %s \n", ip)
			if err != nil {
				return
			}
			return
		}
		// Decrement 1 token for the request
		client.Bucket.Tokens--

		next.ServeHTTP(w, r)
	})
}
