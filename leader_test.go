package main

import (
	"context"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

func TestLeader_DropsToFollowerFromHigherAppendEntryRequestTerm(t *testing.T) {
	timeout := randomDuration(time.Millisecond)
	node := newTestNode(t.Name(), "", timeout, t.Output())
	node.raftState.UpdateState(StateLeader)
	node.raftState.IncrementTerm()

	serverErr := make(chan error, 1)
	go func() {
		node.runLeader(t.Context(), serverErr)
	}()

	replyCh := make(chan RPCReply, 1)

	// send a payload through the network
	// check if it leader mode stops running
	testReqId := NodeId("testNodeId")
	higherTerm := uint64(89)

	sendTimeout, cancel := context.WithTimeout(t.Context(), 1500*time.Millisecond)
	defer cancel()

	select {
	case <-sendTimeout.Done():
		t.Fatalf("leader could not send to network channel while running leader under 1.5s")

	case node.networkCh <- RPCPayload{
		kind: RPCKindAppendEntry,
		payload: AppendEntryRequest{
			Id:   testReqId,
			Term: higherTerm,
		},
		reply: replyCh,
	}:
	}

	resp := <-replyCh
	switch req := resp.payload.(type) {
	case *AppendEntryReply:
		expectedReply := AppendEntryReply{
			Id:     NodeId(t.Name()),
			Term:   higherTerm,
			Result: RaftResultAcked,
		}

		assert.Equal(t, req, &expectedReply)
		actualState := node.raftState.State()
		assert.Equal(t, actualState, StateFollower)
	default:
		t.Fatalf("expected an *AppendEntryReply, got type %T. Reply: %+v\n", req, req)
	}
}
