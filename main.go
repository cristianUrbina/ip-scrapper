package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// generateRandomIP generates a random IPv4 address as a string
func generateRandomIP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d.%d.%d.%d",
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256),
	)
}

// checkAPIAvailability checks common API endpoints on the generated IP
func checkAPIAvailability(ip string, wg *sync.WaitGroup, resultChannel chan<- string) {
	defer wg.Done()

	// List of common API endpoints to check
	endpoints := []string{
		"/api",
		"/public",
		"/v1",
		"/api/v1",
		"/health",
		"/status",
		"/docs",
		"/swagger",
		"/version",
	}

	// Loop over the list of endpoints and check each one
	for _, endpoint := range endpoints {
		url := fmt.Sprintf("http://%s:80%s", ip, endpoint)
		resp, err := http.Get(url)
		if err != nil {
			// If we can't connect to the IP, print an error and move to the next one
			resultChannel <- fmt.Sprintf("Error connecting to %s%s: %v", ip, endpoint, err)
			continue
		}
		defer resp.Body.Close()

		// If the response is OK (200), we have found an API endpoint
		if resp.StatusCode == http.StatusOK {
			resultChannel <- fmt.Sprintf("API found at %s%s (status: %d)", ip, endpoint, resp.StatusCode)
		} else {
			// If not 200, print the status code
			resultChannel <- fmt.Sprintf("No API at %s%s (status: %d)", ip, endpoint, resp.StatusCode)
		}
	}
}

func main() {
	// Define how many IPs to test
	const numIPs = 10000

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Channel to collect results from goroutines
	resultChannel := make(chan string, numIPs*10) // Buffer size to handle multiple results

	// Start testing random IPs
	for i := 0; i < numIPs; i++ {
		randomIP := generateRandomIP()
		wg.Add(1) // Increment the counter for the WaitGroup
		go checkAPIAvailability(randomIP, &wg, resultChannel)
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(resultChannel) // Close the channel once all goroutines are done
	}()

	// Print the results
	for result := range resultChannel {
		fmt.Println(result)
	}
}
