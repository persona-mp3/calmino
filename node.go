package main

import "sync"

type Node struct {
	mu   sync.Mutex
	id   string
	term uint64

	server    Server
	raftState RaftState
	logStore  LogStore
}
