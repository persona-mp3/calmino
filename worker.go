package main

import (
	"context"
	"log/slog"
	"time"
)

// what if the node just has [node.liveCommitIdx] and we give it to the workers? And whenever
// the workers recv a message via replica they then -> liveCommitIdx.Load()? or we just give
// them the [LogStore] because they'll still need to handle snapshotting
type Worker struct {
	// leaderId is the id of the current node
	leaderId NodeId
	// term is the current term of the leader, ideally this should not change
	term  uint64
	store LogStore
	// replicateCh is used when a replication needs to applied across the cluster
	replicateCh chan replicate
	logger      *slog.Logger
}

func NewWorker(leaderId NodeId, term uint64, store LogStore, repCh chan replicate, logger *slog.Logger) *Worker {
	return &Worker{
		leaderId:    leaderId,
		term:        term,
		store:       store,
		replicateCh: repCh,
		logger:      logger,
	}
}

func (w *Worker) Run(leaderCtx context.Context, initialReq AppendEntryRequest, peer RPCConn) {
	d := time.Duration(HeartbeatInterval) * time.Millisecond
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	prevLogEntry := w.store.PreviousEntry()
	commitIdx := w.store.CommitIndex()

	for {
		select {
		case <-leaderCtx.Done():
		case replicateReq := <-w.replicateCh:
			req := replicateReq.req
			reply := AppendEntryReply{}
			if err := peer.Call("Server.AppendEntryRPC", req, &reply); err != nil {
				w.logger.Error("failed to send replication request to", "peer", peer.Id(), "err", err)
				continue
			}

			switch reply.Result {
			case RaftResultAcked, RaftResultLogsOutOfSync:
				select {
				case replicateReq.replicated <- struct{}{}:
					w.logger.Info("sent replication success")
				case <-time.After(ChannelSendTimeout):
					w.logger.Warn("failed to send  acked replication to receiver after", "dur", ChannelSendTimeout)
				}
			default:
				w.logger.Warn("replication was not acked by peer", "reply", reply, "req", req)
			}

			prevLogEntry = w.store.PreviousEntry()
			initialReq.CommitIndex = commitIdx
			initialReq.PreviousLogIndex = prevLogEntry.Index
			initialReq.PreviousLogTerm = prevLogEntry.Term
			commitIdx = w.store.CommitIndex()

		case <-ticker.C:
			reply := AppendEntryReply{}
			// TODO: Add retrials just incase
			if err := peer.Call("Server.AppendEntryRPC", initialReq, &reply); err != nil {
				w.logger.Error("failed to send heertbeat request to", "peer", peer.Id(), "err", err)
				continue
			}

			switch reply.Result {
			case RaftResultAcked:
			case RaftResultLogsOutOfSync:
				w.syncLogsWithPeer(leaderCtx, peer, reply)
			default:
				w.logger.Warn("did not get acked by peer", "id", peer.Id(), "reply", reply)
				return
			}

			ticker.Reset(d)
		}
	}
}

func (w *Worker) syncLogsWithPeer(ctx context.Context, peer RPCConn, reply AppendEntryReply) {
	w.logger.Error("snapshot for syncLogsWithPeer not implemented", "reply", reply)
	panic("snapshot for syncLogsWithPeer not implemented")
}
