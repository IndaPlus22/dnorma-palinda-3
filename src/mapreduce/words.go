package main

//Average runtime: 21.49ms
//Total runtime: 2149 ms
import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"
)

const DataFile = "loremipsum.txt"

// Return the word frequencies of the text argument.
func WordCount(text string) map[string]int {
	cleanText := regexp.MustCompile("[[:punct:]]").ReplaceAllString(text, "")
	words := strings.Fields(strings.ToLower(cleanText))

	numGoroutines := 10
	partialSize := (len(words) + numGoroutines - 1) / numGoroutines

	freqs := make(map[string]int)
	ch := make(chan map[string]int, numGoroutines)

	wg := new(sync.WaitGroup)
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		start := i * partialSize
		stop := (i + 1) * partialSize
		if stop > len(words) {
			stop = len(words)
		}

		go func(slice []string) {
			partialFreq := make(map[string]int)
			for _, word := range slice {
				partialFreq[word]++
			}
			ch <- partialFreq
			wg.Done()
		}(words[start:stop])
	}

	done := false
	go func() {
		wg.Wait()
		close(ch)
		done = true
	}()

	for {
		partialFreq := <-ch
		for word, count := range partialFreq {
			freqs[word] += count
		}
		if(done){
			break
		}
	}
	// This is faster, but sometimes leads to concurrent map read and map write :/
	// go func() {
	// 	for {
	// 		partialFreq := <-ch
	// 		for word, count := range partialFreq {
	// 			freqs[word] += count
	// 		}

	// 	}
	// }()
	return freqs
}

// Benchmark how long it takes to count word frequencies in text numRuns times.
//
// Return the total time elapsed.
func benchmark(text string, numRuns int) int64 {
	start := time.Now()
	for i := 0; i < numRuns; i++ {
		WordCount(text)
	}
	runtimeMillis := time.Since(start).Nanoseconds() / 1e6

	return runtimeMillis
}

// Print the results of a benchmark
func printResults(runtimeMillis int64, numRuns int) {
	fmt.Printf("amount of runs: %d\n", numRuns)
	fmt.Printf("total time: %d ms\n", runtimeMillis)
	average := float64(runtimeMillis) / float64(numRuns)
	fmt.Printf("average time/run: %.2f ms\n", average)
}

func main() {
	// read in DataFile as a string called data
	data, err := ioutil.ReadFile(DataFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v", WordCount(string(data)))

	numRuns := 100
	runtimeMillis := benchmark(string(data), numRuns)
	printResults(runtimeMillis, numRuns)
}
