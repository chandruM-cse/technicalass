package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

type SortRequest struct {
	Data []int json:"data"
}

type SortResponse struct {
	SortedData []int json:"sorted_data"
	TimeTaken int64 json:"time_taken_ms"
}

func main() {
	http.HandleFunc("/process-single", processSingle)
	http.HandleFunc("/process-concurrent", processConcurrent)
	fmt.Println("Server listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func processSingle(w http.ResponseWriter, r *http.Request) {
	var req SortRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	sort.Ints(req.Data)
	endTime := time.Now()

	response := SortResponse{
		SortedData: req.Data,
		TimeTaken:  endTime.Sub(startTime).Milliseconds(),
	}

	json.NewEncoder(w).Encode(response)
}

func processConcurrent(w http.ResponseWriter, r *http.Request) {
	var req SortRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	var wg sync.WaitGroup
	wg.Add(len(req.Data))
	ch := make(chan int)
	for _, item := range req.Data {
		go func(item int) {
			defer wg.Done()
			ch <- item
		}(item)
	}

	sortedData := make([]int, 0, len(req.Data))
	for i := 0; i < len(req.Data); i++ {
		sortedData = append(sortedData, <-ch)
	}
	endTime := time.Now()

	response := SortResponse{
		SortedData: sortedData,
		TimeTaken:  endTime.Sub(startTime).Milliseconds(),
	}

	json.NewEncoder(w).Encode(response)
}