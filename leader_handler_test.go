package main

import (
	"log/slog"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
	"pgregory.net/rapid"
)

func TestLeaderHandlerFromHigherTerm(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(t.Output(), nil))
	handler := NewLeaderHandler(NodeIdNone, logger)
	higherTerm := uint64(8)

	// we can generate requests here
	req := AppendEntryRequest{
		Id:   NodeId(t.Name()),
		Term: higherTerm,
	}

	// generate timeouts
	raftState := NewRaftState(1000 * time.Millisecond)
	logStore := NewLogStore()
	reply, err := handler.AppendEntry(&req, raftState, logStore)
	if err != nil {
		t.Fatalf("leader handler failed: %+v\n", err)
	}

	expectedReply := AppendEntryReply{
		Id:     NodeIdNone,
		Term:   higherTerm,
		Result: RaftResultAcked,
	}

	assert.Equal(t, reply, expectedReply)
	assert.Equal(
		t, raftState.CurrentTerm(), higherTerm, "leader did not update their term",
	)

	assert.Equal(
		t,
		raftState.CurrentLeader(),
		NodeId(t.Name()),
		"leader did not update their leader",
	)

	assert.Equal(t,
		raftState.State(),
		StateFollower,
		"leader did not update their state to folleader")
}

func TestLeaderHandler(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		timeout := time.Duration(rapid.Int32Min(1).Draw(rt, "timeout")) * time.Millisecond
		reqId := rapid.String().Draw(rt, "request-id")
		leaderTerm := rapid.Uint64Range(1, 100).Draw(rt, "leader-term")
		higherTerm := rapid.Uint64Range(1, 300).Draw(rt, "higher-term")

		mockReq := AppendEntryRequest{
			Id:   NodeId(reqId),
			Term: higherTerm,
		}

		logger := slog.New(slog.NewJSONHandler(t.Output(), nil))
		handler := NewLeaderHandler(NodeIdNone, logger)
		raftState := NewRaftState(timeout)
		raftState.state = StateLeader
		raftState.term = leaderTerm
		logStore := NewLogStore()

		reply, err := handler.AppendEntry(&mockReq, raftState, logStore)

		if err != nil {
			t.Fatalf("leader handler failed with error: %s\n", err)
		}

		if leaderTerm == higherTerm {
			expectedReply := AppendEntryReply{
				Id:     NodeIdNone,
				Term:   leaderTerm,
				Result: RaftResultRejectedLeader,
			}
			assert.Equal(
				t, reply, expectedReply,
				"when terms match for a leader it should send a RaftResultRejectedLeader")
		} else if leaderTerm < higherTerm {
			expectedReply := AppendEntryReply{
				Id:     NodeIdNone,
				Term:   higherTerm,
				Result: RaftResultAcked,
			}
			assert.Equal(
				t, reply, expectedReply,
				"when request comes from a higher the leader shoudl send RaftResultAcked")

			assert.Equal(t, raftState.CurrentTerm(), higherTerm)
			assert.Equal(t, raftState.CurrentLeader(), mockReq.Id)
			assert.Equal(t, raftState.State(), StateFollower,
				"when dropping down from leader, state should be a Follower")
		} else {
			expectedReply := AppendEntryReply{
				Id:     NodeIdNone,
				Term:   leaderTerm,
				Result: RaftResultLowerTerm,
			}
			assert.Equal(
				t, reply, expectedReply,
				"when request comes from a lower term the leader should return RaftResultLowerTerm")

			assert.Equal(t, raftState.CurrentTerm(), leaderTerm)
			assert.Equal(t, raftState.State(), StateLeader,
				"node should have remained leader after request came from a lower term")
		}

	})
}
