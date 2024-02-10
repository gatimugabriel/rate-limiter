package bucket

import (
	"fmt"
	"time"
)

func CreateBucket(capacity int, refillInterval time.Duration) *Bucket {
	return &Bucket{
		Capacity:       capacity,
		Tokens:         capacity,
		RefillInterval: time.Second * refillInterval,
		LastRefillTime: time.Now(),
	}
}

func (bucket *Bucket) StartRefillBucket() {
	ticker := time.NewTicker(bucket.RefillInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			bucket.Mu.Lock()
			// refill bucket if not full
			if bucket.Tokens < bucket.Capacity {
				bucket.Tokens++
				fmt.Println("Token added. Total = ", bucket.Tokens)
			}
			bucket.Mu.Unlock()
		}
	}
}
