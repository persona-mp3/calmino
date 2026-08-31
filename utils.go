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
