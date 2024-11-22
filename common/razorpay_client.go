package common

import (
	"fmt"
	"os"
	"sync"
	razorpay "github.com/razorpay/razorpay-go"
)

// Declare a global variable for the Razorpay client
var client *razorpay.Client
var once sync.Once

// RazoryClient initializes the Razorpay client and returns it
func RazoryClient() (*razorpay.Client, error) {
	// Use Once.Do to ensure that the client is initialized only once
	once.Do(func() {
		apiKey := os.Getenv("RAZORPAY_KEY")
		apiSecret := os.Getenv("RAZORPAY_SECRET")

		// Validate the environment variables
		if apiKey == "" || apiSecret == "" {
			client = nil
			return
		}

		// Create the Razorpay client
		client = razorpay.NewClient(apiKey, apiSecret)
	})

	if client == nil {
		return nil, fmt.Errorf("failed to initialize Razorpay client: API key or secret is missing")
	}

	// Return the initialized client
	return client, nil
}
