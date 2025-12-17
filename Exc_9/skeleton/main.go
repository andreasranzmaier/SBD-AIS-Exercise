package main

import (
	"bufio"
	"exc9/mapred"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	// read from file
	lines, err := readMeditations("./res/meditations.txt")
	if err != nil {
		log.Fatalf("failed to read meditations: %v", err)
	}

	// run mapreduce
	var mr mapred.MapReduce
	results := mr.Run(lines)

	// print results sorted by key
	keys := make([]string, 0, len(results))
	for k := range results {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s: %d\n", k, results[k]) // print word counts
	}
}

// readMeditations loads the file into a slice of lines, skipping empty ones.
func readMeditations(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// scan file line by line
	scanner := bufio.NewScanner(f)
	lines := make([]string, 0)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" { // check for empty lines
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}
