package main

import (
	"context"
	"log/slog"
	"time"
)

type replicate struct {
	replicated chan struct{}
	req        AppendEntryRequest
}

type Worker struct {
	id          string
	replicateCh chan replicate
}

var SendChanTimeout = 300 * time.Millisecond

func startReplication(
	ctx context.Context, req AppendEntryRequest, dur time.Duration,
	clusterSize int, workers []*Worker, logger *slog.Logger,
) bool {
	replicated := make(chan struct{}, len(workers)*2) // make sure workers can still send otherwise
	entry := replicate{
		replicated: replicated,
		req:        req,
	}

	for _, worker := range workers {
		go func(w *Worker) {
			ticker := time.NewTicker(SendChanTimeout)
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
			break
		}
	}

	return replicasMade >= quorum
}
