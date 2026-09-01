package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

func (n *Node) runFollower(mainCtx context.Context, serverErrCh chan error) error {
	childCtx, cancel := context.WithCancel(mainCtx)
	defer cancel()

	defer func() {
		n.logger.Info("follower mode exiting")
	}()

	_ = childCtx

	currentState := n.raftState.State()
	if currentState != StateFollower {
		return nil
	}

	n.logger.Info("follower state started successfully")
	electionTimeout := n.raftState.ElectionTimeout()

	ticker := time.NewTicker(electionTimeout)

	fh := NewFollowerHandler(n.id, n.raftState, n.logStore, n.logger)

	for {
		select {
		case <-mainCtx.Done():
			return mainCtx.Err()
		case err := <-serverErrCh:
			return err
		case payload := <-n.networkCh:
			n.logger.Info("recvd payload", slog.Any("payload", payload))
			reply, raftResult := fh.routePayload(payload)
			payload.reply <- reply
			n.logger.Info("payload sent to client: ", slog.Any("payload", reply))
			if raftResult == RaftResultAcked || raftResult == RaftResultLogsOutOfSync {
				ticker.Reset(electionTimeout)
			}

		case <-ticker.C:
			n.logger.Info("did not recv heartbeat from leader, moving to candiadate")
			n.raftState.UpdateState(StateCandidate)
			return nil
		}
	}

}

type followerHandler struct {
	id        string
	raftState *RaftState
	logStore  LogStore
	logger    *slog.Logger
}

func NewFollowerHandler(id string, raftState *RaftState, logStore LogStore, logger *slog.Logger) followerHandler {
	return followerHandler{}
}

func (fh followerHandler) routePayload(payload RPCPayload) (RPCReply, RaftResult) {
	switch req := payload.payload.(type) {
	case AppendEntryRequest:
		return fh.handleAppendEntryRequest(&req)
	case SnapshotRequest:
		return fh.handleSnapshotRequest(&req)
	default:
		panicMsg := fmt.Sprintf("follower has not yet implemented req %v", req)
		panic(panicMsg)
	}
}

// handleAppendEntryRequest checks against the request. If the term from the
// request is the same or higher, it updates the raftStates' term with the
// requests' term and id
func (fh followerHandler) handleAppendEntryRequest(req *AppendEntryRequest) (RPCReply, RaftResult) {
	prevLogEntry := fh.logStore.PreviousEntry()

	reply := AppendEntryReply{
		Id:               NodeId(fh.id),
		Term:             fh.raftState.CurrentTerm(),
		PreviousLogIndex: prevLogEntry.Index,
		PreviousLogTerm:  prevLogEntry.Term,
	}

	raftResult := verifyLeader(req, fh.raftState)

	if raftResult != RaftResultAcked {
		reply.Result = raftResult
		return RPCReply{kind: RPCKindAppendEntry, payload: &reply}, raftResult
	}

	fh.raftState.UpdateTerm(req.Term, req.Id)

	logRaftResult := inspectLogs(req, fh.logStore)

	reply.Result = logRaftResult
	return RPCReply{kind: RPCKindAppendEntry, payload: &reply}, logRaftResult
}

func (fh followerHandler) handleSnapshotRequest(req *SnapshotRequest) (RPCReply, RaftResult) {
	panic("snapshot handler not yet implemented")
}

// inspectLogs checks if the nodes' logs are up to date with the leader by checking
// against the previous log entry of the leader. If none are found, it returns a
// [RaftResultLogsOutOfSync]. If an entry is found at such index but terms don't
// match  it also returns [RaftResultLogsOutOfSync].  Otherwise it returns
// [RaftResultAcked]
func inspectLogs(
	req *AppendEntryRequest, logStore LogStore,
) RaftResult {
	previousLogEntry, err := logStore.LogAt(req.PreviousLogIndex)
	if err != nil {
		return RaftResultLogsOutOfSync
	}
	if previousLogEntry.Term == req.PreviousLogTerm {
		return RaftResultLogsOutOfSync
	}

	return RaftResultAcked
}

// verifyLeader checks if the request came from a valid leader for the current
// term. If the leader is not acked, it returns [RaftResultLowerTerm] or
// [RaftResultRejectedLeader]. Otherwise, [RaftResultAcked]
func verifyLeader(req *AppendEntryRequest, raftState *RaftState) RaftResult {
	currentTerm := raftState.CurrentTerm()
	if req.Term < currentTerm {
		return RaftResultLowerTerm
	}

	if req.Term > currentTerm {
		return RaftResultAcked
	}

	currentLeader := raftState.CurrentLeader()

	switch currentLeader {
	// if we don't have a leader for current term or they match with the our
	// current leader
	case NodeIdNone, req.Id:
		return RaftResultAcked
	default:
		return RaftResultRejectedLeader
	}
}
