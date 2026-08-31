package main

import (
	"context"
	"log/slog"
)

func (n *Node) runCandidate(mainCtx context.Context, serverErrCh chan error) error {
	childCtx, cancel := context.WithCancel(mainCtx)
	defer cancel()
	defer func() {
		n.logger.Info("candidate mode exiting")
	}()

	currentState := n.raftState.State()
	if currentState != StateCandidate {
		return nil
	}

	n.logger.Info("candidate state started successfully")
	_ = childCtx

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
