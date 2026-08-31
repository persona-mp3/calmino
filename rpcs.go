package main

type RPCKind string

const (
	RPCKindSnapshot    RPCKind = "Snapshot"
	RPCKindAppendEntry RPCKind = "AppendEntry"
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
