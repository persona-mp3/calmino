package main

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"
)

type replicate struct {
	// replicated is used to signify an entry that have been successfully replicated
	// on a peer
	replicated chan struct{}
	// req is the append entry to be replicated across the connected peers
	req AppendEntryRequest
}

type Worker struct {
	id        string
	commitIdx *atomic.Uint64
	// replicateCh is used by the leader to start a repication across peers
	replicateCh chan replicate
}

var SendChanTimeout = 300 * time.Millisecond

// startReplication starts a replication across the cluster via the workers.
//
// If a quorum is reached ie a majority of the cluster have replicated the request
// to their logs, it returns true, otherwise false. The [dur] is used to ensure
// that replication has a hard set limit
func startReplication(
	ctx context.Context,
	req AppendEntryRequest,
	dur time.Duration,
	clusterSize int,
	workers []*Worker,
	logger *slog.Logger,
) bool {
	// make sure workers can still send if no receiver
	replicated := make(chan struct{}, len(workers)*2)

	entry := replicate{
		replicated: replicated,
		req:        req,
	}

	for _, worker := range workers {
		go func(w *Worker) {
			ticker := time.NewTicker(SendChanTimeout)
			defer ticker.Stop()
			select {
			case worker.replicateCh <- entry:
			case <-ticker.C:
				logger.Warn("could not send replication entry to worker", "workerId", worker.id)
				return
			}
		}(worker)
	}

	// this node has already made a copy to its logs
	replicasMade := 1
	replicationTimeoutCtx, cancel := context.WithTimeout(ctx, dur)
	quorum := (clusterSize / 2) + 1
	defer cancel()
	for {
		if replicasMade >= quorum {
			return true
		} else if err := replicationTimeoutCtx.Err(); err != nil {
			logger.Info("replication timeout reached")
			break
		}

		select {
		case <-replicated:
			replicasMade += 1
		case <-replicationTimeoutCtx.Done():
		}
	}

	return replicasMade >= quorum
}
