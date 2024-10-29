package bucket

import (
	"testing"
	"time"
)

func TestCreateBucket(t *testing.T) {
	capacity := 10
	refillInterval := 1 * time.Second

	bucket := CreateBucket(capacity, refillInterval)

	if bucket.Capacity != capacity {
		t.Errorf("Expected bucket.Capacity to be %d, got %d", capacity, bucket.Capacity)
	}

	if bucket.Tokens != capacity {
		t.Errorf("Expected bucket.Tokens to be %d, got %d", capacity, bucket.Tokens)
	}

	if bucket.RefillInterval != refillInterval {
		t.Errorf("Expected bucket.RefillInterval to be %s, got %s", refillInterval, bucket.RefillInterval)
	}
}

func TestStartRefillBucket(t *testing.T) {
	capacity := 2
	refillInterval := 500 * time.Millisecond

	// create a bucket with limited tokens for testing
	bucket := CreateBucket(capacity, refillInterval)
	bucket.Tokens = 0 // start with empty bucket

	//start refill process
	go bucket.StartRefillBucket()

	//	wait a little for tokens to be refilled
	time.Sleep(refillInterval + 50*time.Millisecond)

	//check if tokens have been refilled to capacity
	if bucket.Tokens <= 0 {
		t.Errorf("Expected at least 1 token after refill but got %d", bucket.Tokens)
	}
	if bucket.Tokens > capacity {
		t.Errorf("Expected tokens <= capacity, (%d) after refill but got %d", capacity, bucket.Tokens)
	}

}
