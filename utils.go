package main

import (
	"crypto/rand"
	"log"
	"math/big"
	"time"
)

const (
	MaxInterval = 3000
	MinInterval = 2100
)

func randomDuration(d time.Duration) time.Duration {
	limit := big.NewInt(int64(MaxInterval - MinInterval + 1))
	n, err := rand.Int(rand.Reader, limit)
	if err != nil {
		log.Println("warning:: random generator returned 1", n, err)
	}

	actualInterval := n.Int64() + int64(MinInterval)
	return d * time.Duration(actualInterval)
}

func buildIndexedPeerMap(items []string) map[int]map[string][]string {
	result := make(map[int]map[string][]string, len(items))

	for i, currentItem := range items {
		// 1. Allocate exact space for the peer strings
		others := make([]string, 0, len(items)-1)
		others = append(others, items[:i]...)
		others = append(others, items[i+1:]...)

		// 2. Initialize the inner map for this specific index (id)
		result[i] = make(map[string][]string, 1)

		// 3. Assign the item and its peers
		result[i][currentItem] = others
	}

	return result
}
