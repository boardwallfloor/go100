package main

import (
	"fmt"
	"math/rand"
	"slices"
	"sync"
	"time"
)

func unbufferedChan() {
	i := 0
	ch := make(chan struct{})
	go func() {
		i = 1
		<-ch
	}()
	ch <- struct{}{}
	fmt.Println(i)
}

func bufferedChan() {
	i := 0
	ch := make(chan struct{}, 1)
	go func() {
		i = 1
		<-ch
	}()
	ch <- struct{}{}
	fmt.Println(i)
}

func sorting(input []int) []int {
	slices.Sort(input)
	return input
}

func encode(input []int, res chan<- []int, wg *sync.WaitGroup) {
	defer wg.Done() // Ensure the WaitGroup counter is decremented when this goroutine finishes
	for i := 0; i < len(input); i++ {
		if i > 0 {
			time.Sleep(1 * time.Millisecond)
		}
	}
	fmt.Println("Finishing batch", input[0])
	res <- input
}

func parallel2DArray() {
	inputArr := [][]int{
		{1, 2, 2, 2, 3, 3, 3, 4, 5, 7},
		{0, 0, 1, 2, 2, 4, 6, 8, 9, 9},
		{0, 0, 1, 1, 2, 3, 7, 7, 8, 9},
		{0, 1, 3, 5, 6, 7, 7, 8, 8, 9},
		{0, 1, 1, 2, 5, 6, 7, 7, 7, 8},
	}
	// inputArr := generate2dIntArray(4096, 2160)
	// inputArr = randPop2dIntArray(inputArr)
	arrWidth := len(inputArr[0])
	arrHeight := len(inputArr)
	fmt.Printf("Arr width : %d, height : %d, total element :%d\n", arrWidth, arrHeight, arrHeight*arrWidth)
	// chunkLength := 1024
	chunkLength := 7
	maxLength := len(inputArr) * len(inputArr[0])
	chunkCount, leftOvers := getModCount(maxLength, chunkLength)
	fmt.Printf("Estimated chunk %d with leftovers : %d\n", chunkCount, leftOvers)
	chunk := make([]int, chunkLength+1) // +1 to accomodated encoded batch count
	count := 1
	currentChunk := 0

	res := make(chan []int, 10)
	var wg sync.WaitGroup

	for y := range inputArr {
		for x := range inputArr[y] {
			if count <= chunkLength {
				chunk[count] = inputArr[y][x]
				count++
			}
			if count > chunkLength {
				count = 1
				newArr := make([]int, len(chunk))
				copy(newArr, chunk)
				wg.Add(1)
				go encode(newArr, res, &wg)
				currentChunk++
				chunk[0] = currentChunk
			}
			if currentChunk == chunkCount && count > leftOvers {
				newArr := make([]int, len(chunk))
				copy(newArr, chunk)
				wg.Add(1)
				go encode(newArr, res, &wg)
				break
			}
		}
	}
	fmt.Println("Waiting Process")
	go func() {
		wg.Wait()
		close(res)
	}()

	newArr := make([][]int, len(inputArr))
	for i := range newArr {
		newArr[i] = make([]int, len(inputArr[0]))
	}
	for result := range res {
		offset := result[0] * chunkLength
		fmt.Printf("Batch : %d, offset : %d\n", result[0], offset)
		yIndex, xIndex := getModCount(offset, arrWidth)
		// fmt.Printf("Initial x index : %d, y index : %d\n", xIndex, yIndex)
		for count := 0; count < chunkLength; count++ {
			if xIndex >= arrWidth {
				xIndex = 0
				yIndex += 1
			}
			if yIndex >= len(inputArr) {
				break
			}
			// fmt.Printf("X index : %d, Y index : %d\n", xIndex, yIndex)
			newArr[yIndex][xIndex] = result[count+1]
			xIndex++
		}
	}

	// for _, v := range inputArr {
	// 	fmt.Println(v)
	// }
	// fmt.Println("Newly reconstructed")
	// for _, v := range newArr {
	// 	fmt.Println(v)
	// }
	for i := range inputArr {
		for y := range inputArr[i] {
			if inputArr[i][y] != newArr[i][y] {
				fmt.Println("ZONK")
				fmt.Println(i, inputArr[i][y], newArr[i][y])
				break
			}
		}
	}
	fmt.Println("Reconstruction succesful")
}

func getModCount(input, mod int) (int, int) {
	if input%mod == 0 {
		return input / mod, 0
	}
	leftovers := input % mod
	fmt.Printf("Leftovers : %d, Input : %d, Rounded Input : %d\n", leftovers, input, input-leftovers)
	return (input - leftovers) / mod, leftovers
}

func generate2dIntArray(x, y int) [][]int {
	fmt.Println("Generating 2d array")
	arr := make([][]int, y)
	for i := range arr {
		arr[i] = make([]int, x)
	}
	return arr
}

func randPop2dIntArray(arr [][]int) [][]int {
	fmt.Println("Populating 2d array")
	for i := range arr {
		for y := range arr[i] {
			arr[i][y] = rand.Int()
		}
	}
	return arr
}

func main() {
	// stuff()
	parallel2DArray()
}
