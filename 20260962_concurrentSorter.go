// Elliot Buckley (20260962)

package main

// Importing the neccessary imports for the below functions
import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// ConcurrentSorter struct has no fields, serving as a placeholder to attach sorting methods.
type ConcurrentSorter struct{}

// quickSort is a concurrent implementation of the Quick Sort algorithm for sorting an integer slice.
func (cs *ConcurrentSorter) quickSort(data []int) []int {
	if len(data) <= 1 {
		return data // If the data slice is 1 or empty, it's already sorted.
	}

	// Partitioning step of Quick Sort
	left, right := 0, len(data)-1                       // Initialize pointers for the partitioning.
	pivot := rand.Int() % len(data)                     // Choose a random pivot.
	data[pivot], data[right] = data[right], data[pivot] // Swap pivot with the last element.
	for i := range data {
		// If the current element is less than the pivot, swap it to the 'left' part.
		if data[i] < data[right] {
			data[i], data[left] = data[left], data[i]
			left++
		}
	}
	data[left], data[right] = data[right], data[left] // Place the pivot in its correct sorted position.

	// Concurrently sort the partitions.
	leftChan := make(chan []int)
	rightChan := make(chan []int)
	go func() {
		leftChan <- cs.quickSort(data[:left]) // Sort the left partition in a new goroutine.
	}()
	go func() {
		rightChan <- cs.quickSort(data[left+1:]) // Sort the right partition in a new goroutine.
	}()
	return append(<-leftChan, append([]int{data[left]}, <-rightChan...)...) // Combine sorted partitions.
}

// countSort is an implementation of the Counting Sort algorithm, capable of handling negative numbers.
func (cs *ConcurrentSorter) countSort(data []int) []int {
	// Find the range of data for the counting array.
	min, max := data[0], data[0]
	for _, value := range data {
		if value < min {
			min = value // Update min if a smaller value is found.
		}
		if value > max {
			max = value // Update max if a larger value is found.
		}
	}

	offset := 0 // Offset for negative indices.
	if min < 0 {
		offset = -min // Calculate offset based on the minimum value.
	}
	count := make([]int, max-min+1) // Create a counting array including negative values.

	// Populate the counting array.
	for _, value := range data {
		count[value+offset]++ // Increment the count for this number.
	}

	// Construct the sorted array from the counting array.
	z := 0 // Index for the sorted data.
	for i, c := range count {
		for ; c > 0; c-- {
			data[z] = i - offset // Adjust index to get the original number and update the data slice.
			z++
		}
	}
	return data // Return the sorted data.
}

// readCSVFile reads integers from a CSV file and returns them as a slice.
func readCSVFile(filePath string) ([]int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err // Return an error if the file cannot be opened.
	}
	defer file.Close() // Ensure the file is closed after the function exits.

	reader := csv.NewReader(file)
	var data []int // Slice to hold the integers from the CSV file.
	for {
		record, err := reader.Read() // Read a single record from the CSV.
		if err == io.EOF {
			break // If end-of-file is reached, stop reading.
		}
		if err != nil {
			return nil, err // Return an error if the read fails.
		}
		for _, value := range record {
			intValue, err := strconv.Atoi(value) // Convert the string to an integer.
			if err != nil {
				return nil, err // Return an error if conversion fails.
			}
			data = append(data, intValue) // Append the integer to the data slice.
		}
	}
	return data, nil // Return the slice of integers.
}

func main() {

	concurrentSorter := ConcurrentSorter{} // Instantiate the ConcurrentSorter.

	// Define the path to the CSV file containing integers to sort.
	csvFilePath := "20260962_numbers.csv" // Replace with your CSV file path or leave as is if the CSV is in your project folder (VS Code).
	data, err := readCSVFile(csvFilePath) // Read integers from the CSV file.
	if err != nil {
		fmt.Println("Error reading CSV file:", err) // Print an error message if the read fails.
		return
	}

	// DEBUGGING: Un-comment the below print statements to ensure the csv is not sorted prior to starting the sort
	// Print first 10 unsorted elements as a sanity check
	// fmt.Print("\n")
	// fmt.Println("First 10 elements of the csv file NOT sorted: ", data[:10])
	// fmt.Print("\n")

	// Measure Quick Sort Time - Run the sort algorithm 10 times and calculate average time

	var totalQuickElapsed time.Duration // Variable to store the total duration of Quick Sort runs.
	numRuns := 10                       // Number of times to run the sorting algorithms.
	printIterations := true             // Flag to determine whether to print each iteration's results | set to false to disable

	fmt.Printf("\nQUICK SORT ALGORITHM RESULTS:\n\n")

	for i := 0; i < numRuns; i++ {
		dataCopy := make([]int, len(data)) // Create a copy of the data slice for sorting.
		copy(dataCopy, data)               // Copy the original data into the copy.

		start := time.Now()                  // Record the start time of the sort.
		concurrentSorter.quickSort(dataCopy) // Perform the Quick Sort.
		elapsed := time.Since(start)         // Calculate the elapsed time for the sort.
		totalQuickElapsed += elapsed         // Add the elapsed time to the total.

		// Print the sorted data for each run if printIterations is set to true
		if printIterations {

			// DEBUGGING: Un-comment the below line to check that the data is being sorted in each iteration
			// fmt.Printf("Run %d sorted data: %v\n", i+1, concurrentSorter.quickSort(dataCopy)[:10])

			fmt.Printf("Run %d took: %v seconds\n", i+1, elapsed.Seconds())
			fmt.Printf("            %v microseconds\n", elapsed.Microseconds())
			fmt.Printf("            %v nanoseconds\n\n", elapsed.Nanoseconds())
		}
	}

	// Measure Count Sort Time - Run the sort algorithm 10 times and calculate average time

	var totalCountElapsed time.Duration // Variable to store the total duration of Quick Sort runs
	numRuns = 10                        // Number of times to run the sorting algorithms.
	printIterations = true              // Flag to determine whether to print each iteration's results | set to false to disable

	fmt.Printf("COUNT SORT ALGORITHM RESULTS:\n\n")

	for i := 0; i < numRuns; i++ {
		dataCopy := make([]int, len(data)) // Create a copy of the data slice for sorting.
		copy(dataCopy, data)               // Copy the original data into the copy.

		start := time.Now()                  // Record the start time of the sort.
		concurrentSorter.countSort(dataCopy) // Perform the Quick Sort.
		elapsed := time.Since(start)         // Calculate the elapsed time for the sort.
		totalCountElapsed += elapsed         // Add the elapsed time to the total.

		// Print the sorted data for each run if printIterations is set to true
		if printIterations {

			// DEBUGGING: Un-comment the below line to check that the data is being sorted in each iteration
			// fmt.Printf("Run %d sorted data: %v\n", i+1, concurrentSorter.countSort(dataCopy)[:10])

			fmt.Printf("Run %d took: %v seconds\n", i+1, elapsed.Seconds())
			fmt.Printf("            %v microseconds\n", elapsed.Microseconds())
			fmt.Printf("            %v nanoseconds\n\n", elapsed.Nanoseconds())
		}
	}

	// Print the average time taken for the Quick Sort Algorithm to sort the csv file
	averageQuickElapsed := totalQuickElapsed / time.Duration(numRuns)
	fmt.Print("\n")
	fmt.Println("Average Quick sort took: ", averageQuickElapsed.Seconds(), " seconds")
	fmt.Println("                         ", averageQuickElapsed.Microseconds(), " microseconds")
	fmt.Println("                         ", averageQuickElapsed.Nanoseconds(), " nanoseconds")
	fmt.Print("\n")

	// Print the average time taken for the Count Sort Algorithm to sort the csv file
	averageCountElapsed := totalCountElapsed / time.Duration(numRuns)
	fmt.Println("Average Count sort took: ", averageCountElapsed.Seconds(), " seconds")
	fmt.Println("                         ", averageCountElapsed.Microseconds(), " microseconds")
	fmt.Println("                         ", averageCountElapsed.Nanoseconds(), " nanoseconds")
	fmt.Print("\n")
}
