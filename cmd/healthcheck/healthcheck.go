package main

import (
	"context"
	"ec2-benchmarking/lib"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const totalSeconds = 60
const targetUrl = "/liveness"

var requestTimes []int

func makeHTTPRequest(url string, wg *sync.WaitGroup) {
	duration, err := lib.TimeHTTPRequestWaiting(url, wg)
	if err != nil {
		requestTimes = append(requestTimes, duration)
	}
}

func benchmark(url string, n int) {
	var wg sync.WaitGroup
	wg.Add(n)

	limiter := rate.NewLimiter(rate.Every(time.Minute/time.Duration(n)), 1)

	for i := 0; i < n; i++ {
		limiter.Wait(context.Background()) // Block until we are allowed to proceed
		go makeHTTPRequest(url, &wg)
	}

	wg.Wait()
	fmt.Fprintf(os.Stderr, "\n\n")
}

func main() {
	// Check if n is provided as a command-line argument
	if len(os.Args) < 2 {
		fmt.Println("Please provide the value of 'n' as a command-line argument.")
		os.Exit(1)
	}

	// Convert the command-line argument to an integer
	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Invalid value for 'n'. Please provide a valid integer.")
		os.Exit(1)
	}

	// Print current time
	currentTime := time.Now()
	fmt.Printf("Current time: %02d:%02d:%02d\n", currentTime.Hour()+2, currentTime.Minute(), currentTime.Second())
	fmt.Printf("Sending %d requests in %d seconds\n", n, totalSeconds)
	fmt.Println("Please wait...")

	// Perform benchmarking
	benchmark(targetUrl, n)

	// Calculate mean and standard deviation
	sum := 0
	for _, t := range requestTimes {
		sum += t
	}
	mean := float64(sum) / float64(len(requestTimes))

	sumOfSquares := 0.0
	for _, t := range requestTimes {
		sumOfSquares += math.Pow(float64(t)-mean, 2)
	}
	stddev := math.Sqrt(sumOfSquares / float64(len(requestTimes)))

	// Print results
	fmt.Printf("\nResults:\n")
	fmt.Printf("Total requests: %d\n", n)
	fmt.Printf("Valid requests: %d\n", len(requestTimes))
	fmt.Printf("Failed requests: %d\n", n-len(requestTimes))
	fmt.Printf("Time per request: %.2f ± %.2f ms\n", mean, stddev)
}
