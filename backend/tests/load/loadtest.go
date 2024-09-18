package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

var (
	// Seed for randomization, based on current time
	seed       = time.Now().UnixNano()
	// Random source to generate varying values
	randSource = rand.New(rand.NewSource(seed))
)

func main() {
	// Define the duration of the load test (5 minutes)
	duration := 1 * time.Minute
	// Number of concurrent virtual users making requests
	virtualUsers := 100

	// ConstantPacer maintains a consistent request frequency of 60 per second
	pacer := vegeta.ConstantPacer{Freq: 500, Per: time.Second}

	// Define the targets (API endpoints) to test with different HTTP methods
	targets := []*vegeta.Target{
		generateTarget("POST", "http://localhost:8080/api/events/mint", generateMintPayload),
		generateTarget("POST", "http://localhost:8080/api/events/burn", generateBurnPayload),
	}

	// Convert target pointers to values for the Vegeta attack
	targetValues := make([]vegeta.Target, len(targets))
	for i, t := range targets {
		targetValues[i] = *t
	}

	// Create a targeter to handle a fixed list of targets
	targeter := vegeta.NewStaticTargeter(targetValues...)
	// Set up the attacker with the number of workers and a timeout of 30 seconds
	attacker := vegeta.NewAttacker(vegeta.Workers(uint64(virtualUsers)), vegeta.Timeout(30*time.Second))

	// Initialize metrics to track the test's performance
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, pacer, duration, "Load Test") {
		// Log any failed requests for debugging
		if res.Error != "" {
			log.Printf("Request failed: %s\n", res.Error)
		}
		metrics.Add(res)
	}
	metrics.Close()

	// Print summary statistics
	fmt.Printf("Total requests: %d\n", metrics.Requests)
	fmt.Printf("Duration of test: %s\n", metrics.Duration)
	fmt.Printf("Mean requests per second: %.2f\n", metrics.Rate)
	fmt.Printf("99th percentile latency: %s\n", metrics.Latencies.P99)
	fmt.Printf("Average latency: %s\n", metrics.Latencies.Mean)
	fmt.Printf("Minimum latency: %s\n", metrics.Latencies.Min)
	fmt.Printf("Maximum latency: %s\n", metrics.Latencies.Max)
	fmt.Printf("Success ratio: %.2f%%\n", metrics.Success*100)
	fmt.Printf("Average bytes in response: %d bytes\n", int(metrics.BytesIn.Mean))
	fmt.Printf("Total bytes in: %d bytes\n", metrics.BytesIn.Total)
	fmt.Printf("Total bytes out: %d bytes\n", metrics.BytesOut.Total)
}

// Helper function to generate a target with an optional payload
func generateTarget(method, url string, payloadFunc func() []byte) *vegeta.Target {
	target := &vegeta.Target{
		Method: method,
		URL:    url,
	}

	if payloadFunc != nil {
		payload := payloadFunc()
		if len(payload) == 0 {
			log.Printf("Warning: Payload returned empty for %s target\n", method)
			payload = []byte(`{"id":"0x123456","ChainId":"80002","caller_address":"0xC2F20D5c81F5B4450aA9cE62638d0bB01DF1935a","contract_address":"0x1234567890123456789012345678901234567890","block_number":"0x12345678","transaction_hash":"0x1234567890123456789012345678901234567890123456789012345678901234","timestamp":"2023-05-15 12:34:56 MST","amount_from_event":"1000000000000000000","to_from_event":"0x0000000000000000000000000000000000000000"}`)
		}
		target.Body = payload
	}

	return target
}

// Function to generate random payload for mint events
func generateMintPayload() []byte {
	payload := map[string]interface{}{
		"id":                fmt.Sprintf("0x%016x", randSource.Int63()),
		"ChainId":           "80002",
		"caller_address":    "0xC2F20D5c81F5B4450aA9cE62638d0bB01DF1935a",
		"contract_address":  "0x1234567890123456789012345678901234567890",
		"block_number":      randSource.Uint64(),
		"transaction_hash":  fmt.Sprintf("0x%064x", randSource.Uint64()),
		"timestamp":         time.Now().UTC().Format("2006-01-02 15:04:05 MST"),
		"amount_from_event": fmt.Sprintf("%d", randSource.Int63n(1000000000000000000)),
		"to_from_event":     generateRandomAddress(),
	}
	jsonPayload, _ := json.Marshal(payload)
	return jsonPayload
}

// Function to generate random payload for burn events
func generateBurnPayload() []byte {
	payload := map[string]interface{}{
		"id":                fmt.Sprintf("0x%016x", randSource.Int63()),
		"ChainId":           "80002",
		"caller_address":    "0xC2F20D5c81F5B4450aA9cE62638d0bB01DF1935a",
		"contract_address":  "0x1234567890123456789012345678901234567890",
		"block_number":      randSource.Uint64(),
		"transaction_hash":  fmt.Sprintf("0x%064x", randSource.Uint64()),
		"timestamp":         time.Now().UTC().Format("2006-01-02 15:04:05 MST"),
		"amount_from_event": fmt.Sprintf("%d", randSource.Int63n(1000000000000000000)),
		"to_from_event":     "0x0000000000000000000000000000000000000000",
	}
	jsonPayload, _ := json.Marshal(payload)
	return jsonPayload
}

// Helper function to generate a random Ethereum address
func generateRandomAddress() string {
	return fmt.Sprintf("0x%040x", randSource.Uint64())
}

// Pacer to handle dynamic ramp-up and ramp-down of load (not used in this test)
type rampPacer struct {
	start, peak vegeta.Rate
	duration    time.Duration
}

// Function to create a custom ramping pacer for dynamic load tests
func newRampPacer(start, peak float64, duration time.Duration) vegeta.Pacer {
	return &rampPacer{
		start:    vegeta.Rate{Freq: int(start), Per: time.Second},
		peak:     vegeta.Rate{Freq: int(peak), Per: time.Second},
		duration: duration,
	}
}

func (p *rampPacer) Rate(elapsedTime time.Duration) float64 {
	return float64(p.start.Freq)
}

func (p *rampPacer) Pace(elapsedTime time.Duration, elapsedHits uint64) (time.Duration, bool) {
	if elapsedTime >= p.duration {
		return 0, false
	}

	var freq int
	switch {
	case elapsedTime < p.duration/4:
		// Ramp up
		freq = p.start.Freq + int(float64(p.peak.Freq-p.start.Freq)*elapsedTime.Seconds()/(p.duration/4).Seconds())
	case elapsedTime >= p.duration/4 && elapsedTime < p.duration*3/4:
		// Maintain peak
		freq = p.peak.Freq
	default:
		// Ramp down
		remaining := p.duration - elapsedTime
		freq = p.start.Freq + int(float64(p.peak.Freq-p.start.Freq)*remaining.Seconds()/(p.duration/4).Seconds())
	}

	return time.Second / time.Duration(freq), true
}