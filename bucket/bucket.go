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

type Buckets struct {
	Ip string
	Bucket
}
