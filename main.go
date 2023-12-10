package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

// RequestPayload represents the JSON payload structure
type RequestPayload struct {
	ToSort [][]int `json:"to_sort"`
}

// ResponsePayload represents the JSON response structure
type ResponsePayload struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNS       int64   `json:"time_ns"`
}

// SortSequential sorts each sub-array sequentially
func SortSequential(arrays [][]int) ([][]int, int64) {
	startTime := time.Now()

	var sortedArrays [][]int
	for _, arr := range arrays {
		sortedArray := make([]int, len(arr))
		copy(sortedArray, arr)
		sort.Ints(sortedArray)
		sortedArrays = append(sortedArrays, sortedArray)
	}

	endTime := time.Now()
	timeTaken := endTime.Sub(startTime).Nanoseconds()

	return sortedArrays, timeTaken
}

// SortConcurrent sorts each sub-array concurrently
func SortConcurrent(arrays [][]int) ([][]int, int64) {
	startTime := time.Now()

	var wg sync.WaitGroup
	var mutex sync.Mutex
	var sortedArrays [][]int

	for _, arr := range arrays {
		wg.Add(1)
		go func(arr []int) {
			defer wg.Done()

			sortedArray := make([]int, len(arr))
			copy(sortedArray, arr)
			sort.Ints(sortedArray)

			mutex.Lock()
			sortedArrays = append(sortedArrays, sortedArray)
			mutex.Unlock()
		}(arr)
	}

	wg.Wait()

	endTime := time.Now()
	timeTaken := endTime.Sub(startTime).Nanoseconds()

	return sortedArrays, timeTaken
}

func processSingleHandler(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

    fmt.Printf("Person: %v\n", payload)
	sortedArrays, timeTaken := SortSequential(payload.ToSort)

	response := ResponsePayload{
		SortedArrays: sortedArrays,
		TimeNS:       timeTaken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func processConcurrentHandler(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
    
	fmt.Printf("Person: %v\n", payload)
	sortedArrays, timeTaken := SortConcurrent(payload.ToSort)

	response := ResponsePayload{
		SortedArrays: sortedArrays,
		TimeNS:       timeTaken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Define routes
	http.HandleFunc("/process-single", processSingleHandler)
	http.HandleFunc("/process-concurrent", processConcurrentHandler)

	// Start the server
	fmt.Println("Server is running on http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}
