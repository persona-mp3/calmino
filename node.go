package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
)

type Node struct {
	mu sync.Mutex
	// id hold the current id of the node in the cluster
	id string

	listenAddr string

	// peers holds the remote addresses of other peers in the cluster
	peers []string

	// connections hold rpc connections to different nodes
	connections []*RPCConn

	// networkCh is shared with the Server to intercept incoming network RPCs
	networkCh chan RPCPayload

	server    *Server
	raftState *RaftState
	logStore  LogStore

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
	logStore LogStore,
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

	networkCh := make(chan RPCPayload)
	server := NewServer(addr, networkCh, logger)

	return &Node{
		mu:         sync.Mutex{},
		id:         id,
		peers:      peers,
		listenAddr: addr,
		networkCh:  networkCh,
		server:     server,
		raftState:  raftState,
		logStore:   logStore,
		logger:     logger,
	}
}

func (n *Node) Start(mainCtx context.Context) error {
	errCh := make(chan error)
	go func() {
		if err := n.server.Run(mainCtx, "tcp", n.listenAddr); err != nil {
			n.logger.Error("server failed with: ", slog.Any("err", err))
			errCh <- err
		}
	}()

	for {
		switch n.raftState.State() {
		case StateFollower:
			if err := n.runFollower(mainCtx, errCh); err != nil {
				return fmt.Errorf("runFollower err: %w", err)
			}
		case StateCandidate:
			if err := n.runCandidate(mainCtx, errCh); err != nil {
				return fmt.Errorf("runCandidate err: %w", err)
			}
		case StateLeader:
			if err := n.runLeader(mainCtx, errCh); err != nil {
				return fmt.Errorf("runLeader err: %w", err)
			}

		}
	}
}

// Diagnositcs returns string represntation of the nodes current state
func (n *Node) Diagnostics() string {
	return fmt.Sprintf(`Node { id: %s, addr: %s,  peers: %+v, %s`,
		n.id, n.listenAddr, n.peers, n.raftState.String())
}
