package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	// Create a Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Ensure the keyspace notifications are enabled for expired events
	err := rdb.ConfigSet(ctx, "notify-keyspace-events", "Ex").Err()
	if err != nil {
		log.Fatalf("Could not enable keyspace notifications: %v", err)
	}

	// Subscribe to the Redis expired events channel
	pubsub := rdb.PSubscribe(ctx, "__keyevent@0__:expired")
	defer pubsub.Close()

	// Set a key with a TTL to test expiration
	err = rdb.Set(ctx, "mykey", "myvalue", 5*time.Second).Err()
	if err != nil {
		log.Fatalf("Could not set key: %v", err)
	}

	// Start listening for events
	fmt.Println("Waiting for expired events...")
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Fatalf("Error receiving message: %v", err)
		}

		// Handle the expiration event
		fmt.Printf("Key expired: %s\n", msg.Payload)

		// Trigger callback on TTL expire
		handleExpiredKey(msg.Payload)
	}
}

// handleExpiredKey is the callback function triggered on key expiration
func handleExpiredKey(key string) {
	fmt.Printf("Executing callback for expired key: %s\n", key)
	// Implement your callback logic here, e.g., cleaning up resources or notifying other services
}
