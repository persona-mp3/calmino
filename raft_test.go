package main

import (
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

func TestRaft(t *testing.T) {
	initialTimeout := randomDuration(time.Millisecond)
	raftState := NewRaftState(initialTimeout)
	assert.Equal(t, raftState.state, StateFollower, "expected inital state to be follower")

	assert.Equal(t, raftState.electionTimeout, initialTimeout, "expected timeouts to match")

	assert.Equal(t, raftState.CurrentLeader(), NodeIdNone, "expected leader to NodeIdNone")

	nodeId, votedFor := raftState.HasVotedFor(0)
	assert.Equal(t, nodeId, NodeIdNone)
	assert.Equal(t, votedFor, false)

	assert.Equal(t, raftState.CurrentTerm(), uint64(0), "expected current term to be 0")
}
