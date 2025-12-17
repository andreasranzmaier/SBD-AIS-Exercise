package mapred

import (
	"regexp"
	"strings"
	"sync"
)

type MapReduce struct {
}

// Run executes the map and reduce phases concurrently and returns word counts.
func (mr MapReduce) Run(input []string) map[string]int {
	var mapWG sync.WaitGroup
	// channel carries mapper outputs back to the collector
	mapped := make(chan []KeyValue)

	for _, text := range input {
		mapWG.Add(1)
		go func(t string) {
			defer mapWG.Done()
			mapped <- mr.wordCountMapper(t)
		}(text)
	}

	// close channel after all mappers finish
	go func() {
		mapWG.Wait()
		close(mapped)
	}()

	// collect intermediate results grouped by key
	intermediate := make(map[string][]int)
	for kvs := range mapped {
		for _, kv := range kvs {
			intermediate[kv.Key] = append(intermediate[kv.Key], kv.Value)
		}
	}

	// reduce concurrently per key
	results := make(map[string]int)

	// mutex to protect results map
	var reduceWG sync.WaitGroup
	var resMu sync.Mutex

	//
	for key, values := range intermediate {
		reduceWG.Add(1)
		// reducer goroutine
		go func(k string, vals []int) {
			defer reduceWG.Done()
			kv := mr.wordCountReducer(k, vals)
			resMu.Lock()
			results[kv.Key] = kv.Value
			resMu.Unlock()
		}(key, values)
	}

	reduceWG.Wait()
	return results
}

// wordCountMapper converts a string into lowercase and returns a KeyValue per word
func (mr MapReduce) wordCountMapper(text string) []KeyValue {
	wordRegex := regexp.MustCompile(`[A-Za-z]+`)
	words := wordRegex.FindAllString(text, -1)
	result := make([]KeyValue, 0, len(words))
	for _, w := range words {
		result = append(result, KeyValue{Key: strings.ToLower(w), Value: 1}) // normalize to lowercase
	}
	return result
}

// wordCountReducer sums up the counts for a given word
func (mr MapReduce) wordCountReducer(key string, values []int) KeyValue {
	sum := 0
	for _, v := range values {
		sum += v
	}
	return KeyValue{Key: key, Value: sum}
}
