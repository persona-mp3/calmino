package main

import (
	"io"
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

func NewNode(id, addr string,
	peers []string,
	logStore *LogStore,
	raftState *RaftState,
	server *Server,
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

	logger.Info("started logging")
	return &Node{
		mu:        sync.Mutex{},
		id:        id,
		server:    server,
		raftState: raftState,
		logStore:  logStore,
		logger:    logger,
	}
}
