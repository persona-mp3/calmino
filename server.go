package main

import "log/slog"

type RPC struct{}

type Server struct {
	network chan RPC
	log     *slog.Logger
}
