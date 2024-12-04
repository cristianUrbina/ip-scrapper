package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

func generateRandomIP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d.%d.%d.%d",
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256),
	)
}

func checkAPIAvailability(ip string, wg *sync.WaitGroup, resultChannel chan<- string) {
	defer wg.Done()
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
	for _, endpoint := range endpoints {
		url := fmt.Sprintf("http://%s:80%s", ip, endpoint)
		resp, err := http.Get(url)
		if err != nil {
			resultChannel <- fmt.Sprintf("Error connecting to %s%s: %v", ip, endpoint, err)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			resultChannel <- fmt.Sprintf("API found at %s%s (status: %d)", ip, endpoint, resp.StatusCode)
		} else {
			resultChannel <- fmt.Sprintf("No API at %s%s (status: %d)", ip, endpoint, resp.StatusCode)
		}
	}
}

func main() {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	const numIPs = 10000
	var wg sync.WaitGroup
	resultChannel := make(chan string, numIPs*10)
	for i := 0; i < numIPs; i++ {
		randomIP := generateRandomIP()
		wg.Add(1)
		go checkAPIAvailability(randomIP, &wg, resultChannel)
	}
	go func() {
		wg.Wait()
		close(resultChannel)
	}()
	for result := range resultChannel {
		log.Println(result)
	}
}
