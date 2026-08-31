package main

import "log/slog"

type RPC struct{}

type Server struct {
	network chan RPC
	logger  *slog.Logger
}

func NewServer(addr string, networkCh chan RPC) *Server {
	return &Server{network: networkCh}
}
