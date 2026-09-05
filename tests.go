package main

import (
	"io"
	"time"
)

func newTestNode(id, addr string, timeout time.Duration, out io.Writer) *Node {
	store := NewLogStore()
	state := NewRaftState(timeout)
	node := NewNode(id, addr, []string{}, store, state, Configuration{
		Out: out,
	})
	return node
}
