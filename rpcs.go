package main

type RPCKind string

const (
	RPCKindSnapshot    RPCKind = "Snapshot"
	RPCKindAppendEntry RPCKind = "AppendEntry"
	RPCKindVote        RPCKind = "Vote"
)

type RPCReply struct {
	kind    RPCKind
	payload any
}

type RPCPayload struct {
	kind    RPCKind
	payload any
	reply   chan RPCReply
}

type SnapshotRequest struct {
	Id               NodeId
	Term             uint64
	Result           RaftResult
	PreviousLogIndex uint64
	PreviousLogTerm  uint64
}

type SnapshotReply struct {
	Id          NodeId
	Term        uint64
	Result      RaftResult
	Snapshot    []Log
	CommitIndex uint64
}

type AppendEntryRequest struct {
	Id               NodeId
	Term             uint64
	Result           RaftResult
	PreviousLogIndex uint64
	PreviousLogTerm  uint64
}

type AppendEntryReply struct {
	Id               NodeId
	Term             uint64
	Result           RaftResult
	PreviousLogIndex uint64
	PreviousLogTerm  uint64
}

// VoteRequest should have a Result of RaftResultVoteRequest
type VoteRequest struct {
	Id               NodeId
	Term             uint64
	Result           VoteResult
	Message          string
	PreviousLogIndex uint64
	PreviousLogTerm  uint64
}

// VoteReply should have a Result of RaftResultVoteGranted or RaftResultVoteDenied
// and then a more detailed reason inside message like a string. Or maybe we could
// just set Result to be []RaftResult? or add extend it to have a Reasons field?
type VoteReply struct {
	Id               NodeId
	Term             uint64
	Result           VoteResult
	Message          string
	PreviousLogIndex uint64
	PreviousLogTerm  uint64
}
