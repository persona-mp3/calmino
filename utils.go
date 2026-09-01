package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	_ "net/http/pprof"
	"net/rpc"
	"os"
	"strconv"
	"time"
)

func randomDuration(d time.Duration) time.Duration {
	limit := big.NewInt(int64(ELECTION_INTERVAL_MAX - ELECTION_INTERVAL_MIN + 1))
	n, err := rand.Int(rand.Reader, limit)
	if err != nil {
		log.Println("warning:: random generator returned 1", n, err)
	}

	actualInterval := n.Int64() + int64(ELECTION_INTERVAL_MIN)
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

func createFileWithName(name string) io.Writer {
	f, err := os.Create(name)
	if err != nil {
		f = os.Stdout
		log.Printf(
			"[warn] could not create file: %s, reason: %s. using stdout\n",
			name,
			err)
	}

	return f
}

func startPprofServer(addr string) error {
	if err := http.ListenAndServe(addr, nil); err != nil {
		return fmt.Errorf("failed to start pprof-server: %w", err)
	}
	return nil
}

func connectToPeers(network string, addrs []string) []*RPCPeer {
	rpcPeers := []*RPCPeer{}
	for idx, addr := range addrs {
		conn, err := rpc.Dial(network, addr)
		if err != nil {
			log.Printf("[err] could not dial %s. reason: %s\n", addr, err)
			continue
		}
		id := strconv.Itoa(idx)
		rpcPeers = append(rpcPeers, &RPCPeer{id: id, addr: addr, conn: conn})
	}
	return rpcPeers
}
