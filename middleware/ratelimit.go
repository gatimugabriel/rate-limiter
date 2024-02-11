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

			for ip, client := range clients {
				if time.Since(client.LastSeen) > 1*time.Minute {
					fmt.Println("\t Killed client...", ip)
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

		fmt.Println("clients", clients)

		mu.Lock()
		if _, exists := clients[ip]; !exists {
			fmt.Println("clients[ip]", clients[ip])
			fmt.Println(exists)
			fmt.Println("New client initiated\n ")
			// create client bucket
			clients[ip] = &Client{
				Bucket: bucket.CreateBucket(10, 1),
			}

			fmt.Println("created client", clients[ip])
			fmt.Println("clients", &clients)

			// ---@routine -> start the bucket refill ()
			//clients[ip].Bucket.StartRefillBucket()
		} else {
			fmt.Println("Client retrieved, process continue...\n ")
		}

		//fmt.Println("\t |---------- Clients(After) ----------|\n\t")
		//for ipAddr, client := range clients {
		//	fmt.Printf("\t IP: %s, Bucket: %+v\n", ipAddr, client.Bucket)
		//}

		clients[ip].LastSeen = time.Now()

		//if clients[ip].Bucket.Tokens == 0 {
		//	mu.Unlock()
		//	w.WriteHeader(http.StatusTooManyRequests)
		//	_, err := fmt.Fprintf(w, "Too many requests for IP: %s", ip)
		//	if err != nil {
		//		return
		//	}
		//	return
		//}
		
		// Decrement token for allowed request
		//clients[ip].Bucket.Tokens--
		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
