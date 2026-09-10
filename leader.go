package main

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

func (n *Node) runLeader(mainCtx context.Context, serverErrCh chan error) error {
	n.logger.Info("leader state started successfully")
	handler := NewLeaderHandler(NodeId(n.id), n.logger)

	// Assumes the candidate left connections open and persistent and start workers on them
	currentTerm := n.raftState.CurrentTerm()
	previousLogEntry := n.logStore.PreviousEntry()
	commitIdx := n.logStore.CommitIndex()
	workerWg := sync.WaitGroup{}

	leaderCtx, cancel := context.WithCancel(mainCtx)
	defer func() {
		cancel()
		n.raftState.UpdateState(StateFollower)
		for _, peer := range n.rpcConnections {
			if err := peer.Close(); err != nil {
				n.logger.Error("could not close connection", "id", peer.Id(), "err", err)
			}
		}
	}()

	initialReq := AppendEntryRequest{
		Id:               NodeId(n.id),
		Term:             currentTerm,
		PreviousLogIndex: previousLogEntry.Index,
		PreviousLogTerm:  previousLogEntry.Term,
		CommitIndex:      commitIdx,
	}
	for _, peer := range n.rpcConnections {
		replicateCh := make(chan replicate)
		worker := NewWorker(NodeId(n.id), currentTerm, n.logStore, replicateCh, n.logger)
		workerWg.Go(func() { worker.Run(leaderCtx, initialReq, peer) })

	}

	workersReturned := make(chan struct{})
	go func() {
		// this node only steps down from a leader, the whole cluster refused to ack
		// them through the workers. This is to ensure that a single Follower in the
		// cluster who perchance has a higher term does not negate the consensus if
		// other Followers in the clusters are happy with this leader.
		workerWg.Wait()
		close(workersReturned)
	}()

	for {
		currentState := n.raftState.State()
		if currentState != StateLeader {
			n.logger.Info(
				"leader state changed. now dropping to", slog.Any("state", currentState),
			)
			break
		}

		select {
		case <-mainCtx.Done():
			return mainCtx.Err()
		case err := <-serverErrCh:
			return err
		case rpcPayload := <-n.networkCh:
			n.logger.Info("received rpcPayload", slog.Any("payload", rpcPayload))
			reply, _, err := router(rpcPayload, n.raftState, n.logStore, handler)
			if err != nil {
				// TODO: Proper error handling
				panic(err)
			}

			rpcPayload.reply <- reply
		case <-workersReturned:
			n.logger.Warn("all workers have returned, exiting leader state")
			return nil
		}
	}

	return nil
}

func NewLeaderHandler(id NodeId, logger *slog.Logger) Handler {
	return &leaderHandler{id: id, logger: logger}
}

func router(
	payload RPCPayload,
	raftState *RaftState,
	logStore LogStore,
	handler Handler,
) (RPCReply, RaftResult, error) {
	switch req := payload.payload.(type) {
	case AppendEntryRequest:
		reply, err := handler.AppendEntry(&req, raftState, logStore)
		return RPCReply{kind: RPCKindAppendEntry, payload: &reply}, reply.Result, err
	case VoteRequest:
		reply, err := handler.Vote(&req, raftState, logStore)
		return RPCReply{kind: RPCKindAppendEntry, payload: &reply}, reply.Result, err
	case SnapshotRequest:
		reply, err := handler.Snapshot(&req, raftState, logStore)
		return RPCReply{kind: RPCKindAppendEntry, payload: &reply}, reply.Result, err
	}
	errMsg := fmt.Sprintf(`unexpected payload to router as leader, kind: %+v payload: %+v `, payload.kind, payload.payload)
	return RPCReply{}, RaftResultUnknownUnhandled, fmt.Errorf("%s", errMsg)
}

func (lh *leaderHandler) AppendEntry(
	req *AppendEntryRequest, raftState *RaftState, logStore LogStore,
) (AppendEntryReply, error) {
	currentTerm := raftState.CurrentTerm()
	prevLog := logStore.PreviousEntry()

	if req.Term < currentTerm {
		return AppendEntryReply{
			Id:               lh.id,
			Term:             currentTerm,
			Result:           RaftResultLowerTerm,
			PreviousLogIndex: prevLog.Index,
			PreviousLogTerm:  prevLog.Term,
		}, nil
	}

	if req.Term == currentTerm {
		lh.logger.Warn("recvd appendEntry with the same term from another node")
		return AppendEntryReply{
			Id:               lh.id,
			Term:             currentTerm,
			Result:           RaftResultRejectedLeader,
			PreviousLogIndex: prevLog.Index,
			PreviousLogTerm:  prevLog.Term,
		}, nil
	}

	raftState.UpdateTerm(req.Term, req.Id)
	raftState.UpdateState(StateFollower)
	lh.logger.Info(
		"dropping from leader to follower due to higher term from request",
		slog.Uint64("currentTerm", currentTerm), slog.Any("request", req),
	)
	return AppendEntryReply{
		Id:               lh.id,
		Term:             req.Term,
		Result:           RaftResultAcked,
		PreviousLogIndex: prevLog.Index,
		PreviousLogTerm:  prevLog.Term,
	}, nil
}

func (lh *leaderHandler) Vote(
	req *VoteRequest, raftState *RaftState, logStore LogStore,
) (VoteReply, error) {
	currentTerm := raftState.CurrentTerm()
	prevLogEntry := logStore.PreviousEntry()
	if req.Term <= currentTerm {
		return VoteReply{
			Id:               lh.id,
			Term:             currentTerm,
			Result:           RaftResultVoteDenied,
			PreviousLogIndex: prevLogEntry.Index,
			PreviousLogTerm:  prevLogEntry.Term,
		}, nil
	}

	logComplete := logCompletion(
		prevLogEntry, req.PreviousLogIndex, req.PreviousLogTerm,
	)

	if !logComplete {
		lh.logger.Info("logs from candidate are incomplete",
			slog.Any("prevLogEntry", prevLogEntry),
			slog.Any("reqPrevLogIndex", req.PreviousLogIndex),
			slog.Any("reqPrevLogTerm", req.PreviousLogTerm),
		)
		return VoteReply{
			Id:               lh.id,
			Term:             currentTerm,
			Result:           RaftResultVoteDenied,
			Message:          "incomplete logs",
			PreviousLogIndex: prevLogEntry.Index,
			PreviousLogTerm:  prevLogEntry.Term,
		}, nil
	}
	raftState.GrantVoteTo(req.Term, req.Id)
	raftState.UpdateState(StateFollower)
	lh.logger.Info("logs from candidate are complete",
		slog.Any("prevLogEntry", prevLogEntry),
		slog.Any("reqPrevLogIndex", req.PreviousLogIndex),
		slog.Any("reqPrevLogTerm", req.PreviousLogTerm),
	)

	lh.logger.Info(
		"recvd voteRPC from a higher term than leader, yeilding to them",
		slog.Uint64("currentTerm", currentTerm), slog.Any("reqTerm", req),
	)

	return VoteReply{
		Id:               lh.id,
		Term:             req.Term,
		Result:           RaftResultVoteGranted,
		PreviousLogIndex: prevLogEntry.Index,
		PreviousLogTerm:  prevLogEntry.Term,
	}, nil
}

func (lh *leaderHandler) Snapshot(
	req *SnapshotRequest, raftState *RaftState, logStore LogStore,
) (SnapshotReply, error) {
	panic("leader-snapshot request not implemented")
}
