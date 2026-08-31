package main

import (
	"context"
	"log/slog"
)

func (n *Node) runLeader(mainCtx context.Context, serverErrCh chan error) error {
	childCtx, cancel := context.WithCancel(mainCtx)
	defer cancel()

	_ = childCtx
	currentState := n.raftState.State()
	if currentState != StateLeader {
		return nil
	}

	n.logger.Info("leader state started successfully")

	for {
		select {
		case <-mainCtx.Done():
			return mainCtx.Err()
		case err := <-serverErrCh:
			return err
		case payload := <-n.networkCh:
			n.logger.Info("recvd payload", slog.Any("payload", payload))
		}
	}
}
