package middleware

import (
	"fmt"
	"github.com/gabrielgatimu/rate-limiter/bucket"
	"log"
	"net/http"
	"sync"
	"time"
)

func RateLimit(next http.Handler) http.Handler {
	type Client struct {
		Bucket   *bucket.Bucket
		LastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*Client)
	)

	// ---@routine -> kill old clients
	go func() {
		for {
			time.Sleep(1 * time.Second)
			mu.Lock()

			fmt.Println("Killer in action...")
			for ip, client := range clients {
				if time.Since(client.LastSeen) > 1*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, ok := r.Context().Value("ip").(string)
		if !ok {
			log.Fatal("Missing IP")
		}

		fmt.Println("\t |---------- Clients(Before) ----------|\n  \t")
		for ip, client := range clients {
			fmt.Printf("\t IP: %s, Bucket: %+v\n", ip, client.Bucket)
		}

		mu.Lock()
		if _, exists := clients[ip]; !exists {
			fmt.Println("New client initiated")
			// create client bucket
			clients[ip] = &Client{
				Bucket: bucket.CreateBucket(10, 1),
			}

			// ---@routine -> start the bucket refill ()
			go clients[ip].Bucket.StartRefillBucket()
		} else {
			fmt.Println("Client retrieved, process continue...")
		}

		fmt.Println("\t |---------- Clients(After) ----------|\n\t")
		for ip, client := range clients {
			fmt.Printf("\t IP: %s, Bucket: %+v\n", ip, client.Bucket)
		}

		clients[ip].LastSeen = time.Now()

		if clients[ip].Bucket.Tokens == 0 {
			mu.Unlock()
			w.WriteHeader(http.StatusTooManyRequests)
			_, err := fmt.Fprintf(w, "Too many requests for IP: %s", ip)
			if err != nil {
				return
			}
			return
		}
		mu.Unlock()

		// Decrement token for allowed request
		clients[ip].Bucket.Tokens--

		next.ServeHTTP(w, r)
	})
}
