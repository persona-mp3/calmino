package main

import "log/slog"

type Handler interface {
	AppendEntry(
		req *AppendEntryRequest,
		raftState *RaftState,
		logStore LogStore,
	) (AppendEntryReply, error)

	Vote(
		req *VoteRequest,
		raftState *RaftState,
		logStore LogStore,
	) (VoteReply, error)

	Snapshot(
		req *SnapshotRequest,
		raftState *RaftState,
		logStore LogStore,
	) (SnapshotReply, error)
}

type leaderHandler struct {
	id     NodeId
	logger *slog.Logger
}

type candidateHandler struct {
	id     NodeId
	logger *slog.Logger
}

func NewCandidateHandler(
	id NodeId,
	logger *slog.Logger,
) Handler {
	return &candidateHandler{
		id:     id,
		logger: logger,
	}
}

func (ch *candidateHandler) AppendEntry(
	req *AppendEntryRequest,
	raftState *RaftState,
	logStore LogStore,
) (AppendEntryReply, error) {
	currentTerm := raftState.CurrentTerm()
	prevLogEntry := logStore.PreviousEntry()
	if req.Term < currentTerm {
		return AppendEntryReply{
			Id:               ch.id,
			Term:             currentTerm,
			Result:           RaftResultLowerTerm,
			PreviousLogIndex: prevLogEntry.Index,
			PreviousLogTerm:  prevLogEntry.Term,
		}, nil
	}

	ch.logger.Info(
		"recvd appendEntry term  higher than candidates, yeilding to them",
		slog.Uint64("currentTerm", currentTerm),
		slog.Uint64("reqTerm", req.Term),
	)

	return AppendEntryReply{
		Id:               ch.id,
		Term:             req.Term,
		Result:           RaftResultAcked,
		PreviousLogIndex: prevLogEntry.Index,
		PreviousLogTerm:  prevLogEntry.Term,
	}, nil
}

func (ch *candidateHandler) Vote(
	req *VoteRequest,
	raftState *RaftState,
	logStore LogStore,
) (VoteReply, error) {
	currentTerm := raftState.CurrentTerm()
	prevLogEntry := logStore.PreviousEntry()
	/// TODO: Impl log checking
	if req.Term <= currentTerm {
		return VoteReply{
			Id:               ch.id,
			Term:             currentTerm,
			Result:           RaftResultVoteDenied,
			PreviousLogIndex: prevLogEntry.Index,
			PreviousLogTerm:  prevLogEntry.Term,
		}, nil
	}

	raftState.GrantVoteTo(req.Term, req.Id)
	ch.logger.Info(
		"recvd voteRPC from a higher term than candidates, yeilding to them",
		slog.Uint64("currentTerm", currentTerm),
		slog.Uint64("reqTerm", req.Term),
	)
	return VoteReply{
		Id:               ch.id,
		Term:             req.Term,
		Result:           RaftResultVoteGranted,
		PreviousLogIndex: prevLogEntry.Index,
		PreviousLogTerm:  prevLogEntry.Term,
	}, nil
}

func (ch *candidateHandler) Snapshot(
	req *SnapshotRequest,
	raftState *RaftState,
	logStore LogStore,
) (SnapshotReply, error) {
	panic("candidate not impl snapshot")
}
