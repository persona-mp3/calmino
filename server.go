package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/rpc"
)

type Server struct {
	network chan RPCPayload
	logger  *slog.Logger
}

func NewServer(addr string, networkCh chan RPCPayload, logger *slog.Logger) *Server {
	return &Server{network: networkCh, logger: logger}
}

func (s *Server) Run(ctx context.Context, network, addr string) error {
	rpcHandler := rpc.NewServer()
	if err := rpcHandler.Register(s); err != nil {
		return fmt.Errorf("could not register rpc server: %w", err)
	}
	ln, err := net.Listen(network, addr)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		if err := ln.Close(); err != nil {
			s.logger.Info("could not close listener:", slog.Any("err", err))
		}
		s.logger.Info("shutting down server")
	}()

	log.Println("rpcServer started at", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				s.logger.Debug("listener  has been closed")
				return nil
			}
			s.logger.Error("could not accept connection: ", slog.Any("err", err))
			continue
		}

		go rpcHandler.ServeConn(conn)

	}
}
