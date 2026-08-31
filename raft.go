package main

import (
	"sync"
	"time"
)

type State string
type NodeId string

const (
	StateLeader    State = "Leader"
	StateFollower  State = "Follower"
	StateCandidate State = "Candidate"
)

const (
	NodeIdNone NodeId = ""
)

type RaftState struct {
	mu sync.Mutex
	// stores current term for the node
	term uint64
	// votedFor stores the various terms and who this node has voted for
	votedFor map[uint64]NodeId
	// leader stores the leader for that was acknowledged by the node
	leader NodeId
	// electionTimeout stores the electionTimeout for the node
	electionTimeout time.Duration

	state State
}

func NewRaftState(initialTimeout time.Duration) *RaftState {
	voteHistory := make(map[uint64]NodeId)
	return &RaftState{
		mu:              sync.Mutex{},
		term:            0,
		leader:          NodeIdNone,
		votedFor:        voteHistory,
		electionTimeout: initialTimeout,
		state:           StateFollower,
	}
}

// IncrementTerm increments the term for this node and returns the new term
func (rs *RaftState) IncrementTerm() uint64 {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.term++
	return rs.term
}

// CurrentTerm returns the current term
func (rs *RaftState) CurrentTerm() uint64 {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return rs.term
}

func (rs *RaftState) UpdateTerm(term uint64, leader NodeId) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.term = term
	rs.leader = leader
}

// GrantVoteTo grants the term and vote to a candidate
func (rs *RaftState) GrantVoteTo(term uint64, candidateId NodeId) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.votedFor[term] = candidateId
}

// CurrentLeader returns the current leader for this term
func (rs *RaftState) CurrentLeader() NodeId {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return rs.leader
}

// HasVotedFor checks if this node has voted for term. It returns true of a vote
// for this term has been found
func (rs *RaftState) HasVotedFor(term uint64) (NodeId, bool) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	nodeId, ok := rs.votedFor[term]
	return nodeId, ok
}

func (rs *RaftState) ClearLeader() {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.leader = NodeIdNone
}
