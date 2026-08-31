package main

import (
	"log/slog"
	"sync"
)

type Node struct {
	mu sync.Mutex
	// id hold the current id of the node in the cluster
	id string

	// peers holds the remote addresses of other peers in the cluster
	peers []string

	// connections hold rpc connections to different nodes
	connections []*RPCConn

	server    Server
	raftState RaftState
	logStore  LogStore

	logger *slog.Logger
}
