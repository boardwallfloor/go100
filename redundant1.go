package main

import (
	"fmt"
	"sync"
	"time"
)

func processChunk(chunk []int, res chan<- []int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := range chunk {
		chunk[i] *= 10
		time.Sleep(1 * time.Second) // simulate processing delay
	}
	res <- chunk
}

// ChunkArray divides the array into smaller chunks of specified size.
func chunkArray(array []int, chunkSize int) [][]int {
	var chunks [][]int
	for i := 0; i < len(array); i += chunkSize {
		end := i + chunkSize
		if end > len(array) {
			end = len(array)
		}
		chunks = append(chunks, array[i:end])
	}
	return chunks
}

func stuff() {
	inputArray := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	chunkSize := 1

	// Chunk the array into smaller sub-arrays
	chunks := chunkArray(inputArray, chunkSize)

	fmt.Println(len(chunks))
	res := make(chan []int, len(chunks)) // Channel to collect results
	var wg sync.WaitGroup

	// Process each chunk in a separate goroutine
	for _, chunk := range chunks {
		wg.Add(1)
		go processChunk(chunk, res, &wg)
	}

	// Start a goroutine to process results as they come in
	go func() {
		wg.Wait()
		close(res) // Close the channel once all processing is done
	}()

	// Collect and process results from the result channel
	count := 0
	for result := range res {
		fmt.Println("Processed chunk:", result)
		fmt.Println(count)
		count++
	}

	fmt.Println("All chunks processed.")
}
