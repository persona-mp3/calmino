package main

import (
	"io"
	"log/slog"
	"sync"
	"time"
)

type Node struct {
	mu sync.Mutex
	// id hold the current id of the node in the cluster
	id string

	// peers holds the remote addresses of other peers in the cluster
	peers []string

	// connections hold rpc connections to different nodes
	connections []*RPCConn

	// networkCh is shared with the Server to intercept incoming network RPCs
	networkCh chan RPC

	server    *Server
	raftState *RaftState
	logStore  *LogStore

	logger        *slog.Logger
	configuration Configuration
}

type Configuration struct {
	LogLevel  slog.Level
	LogFormat int
	Addrs     []string
	Out       io.Writer
}

func NewNode(
	id,
	addr string,
	peers []string,
	logStore *LogStore,
	raftState *RaftState,
	config Configuration,
) *Node {
	var logger *slog.Logger

	if config.LogFormat == 0 {
		logger = slog.New(
			slog.NewTextHandler(
				config.Out, &slog.HandlerOptions{Level: config.LogLevel},
			))
	} else {
		logger = slog.New(
			slog.NewJSONHandler(
				config.Out, &slog.HandlerOptions{Level: config.LogLevel},
			))
	}

	networkCh := make(chan RPC)
	server := NewServer(networkCh)

	return &Node{
		mu:        sync.Mutex{},
		id:        id,
		networkCh: networkCh,
		server:    server,
		raftState: raftState,
		logStore:  logStore,
		logger:    logger,
	}
}

func (n *Node) Run() error {
	electionTimeout := n.raftState.ElectionTimeout()
	<-time.After(electionTimeout)
	// TODO: startServer
	n.logger.Info("node recvd no hearbeat, transition to candidate")
	return nil
}
