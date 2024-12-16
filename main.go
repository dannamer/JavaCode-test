package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	apiURL = "http://localhost:8080/api/v1/wallet"
	requestsPerSec = 1000
)

type WalletRequest struct {
	WalletID      string  `json:"walletId"`
	OperationType string  `json:"operationType"`
	Amount        float64 `json:"amount"`
}

func sendRequest(wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()

	walletID := "8defb3ed-96be-4e98-857f-d0ff09e5e56d"
	operationType := "DEPOSIT"
	amount := 1000.0

	requestBody := WalletRequest{
		WalletID:      walletID,
		OperationType: operationType,
		Amount:        amount,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		errChan <- fmt.Errorf("failed to marshal request body: %v", err)
		return
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		errChan <- fmt.Errorf("failed to send request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 && resp.StatusCode < 600 {
		errChan <- fmt.Errorf("received 50X error: %d", resp.StatusCode)
	}
}

func main() {
	var wg sync.WaitGroup
	errChan := make(chan error, requestsPerSec)

	startTime := time.Now()
	var count int

	for i := 0; i < requestsPerSec; i++ {
		wg.Add(1)
		go sendRequest(&wg, errChan)
		count++
	}

	wg.Wait()

	close(errChan)
	for err := range errChan {
		if err != nil {
			log.Println("Error:", err)
		}
	}

	elapsedTime := time.Since(startTime)
	if len(errChan) > 0 {
		fmt.Printf("Test finished with %d errors in %v\n", len(errChan), elapsedTime)
	} else {
		fmt.Printf("Test completed successfully in %v. No 50X errors. Total requests: %d\n", elapsedTime, count)
	}
}
