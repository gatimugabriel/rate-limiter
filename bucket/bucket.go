package bucket

import (
	"sync"
	"time"
)

type Bucket struct {
	Capacity       int
	Tokens         int
	RefillInterval time.Duration
	LastRefillTime time.Time
	Mu             sync.Mutex		
}

type Client struct {
	Bucket
	LastSeen time.Time
}
