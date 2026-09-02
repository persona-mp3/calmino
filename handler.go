package main

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

type leaderHandler struct{}
type candidateHandler struct{}
