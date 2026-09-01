// QUESTION: What if just kept all our connections from startup? There would
// be no need for redials and such, but is holding a connection open for long
// worth it? we might not know they've disconnected till we want to send to them?
//
//
//
// Extra docs and references
// what if we just give collectOtherVotes to do all the voting logic and
// it just tells this loop where it won as RaftResultEelected || RaftResultLostElection

// 5.4.1 Election Restriction
// Raft determines which of two logs is more up-to-date
// by comparing the index and term of the last entries in the
// logs. If the logs have last entries with different terms, then
// the log with the later term is more up-to-date. If the logs
// end with the same term, then whichever log is longer is
// more up-to-date.

// While waiting for votes, a candidate may receive an
// AppendEntries RPC from another server claiming to be
// leader. If the leader’s term (included in its RPC) is at least
// as large as the candidate’s current term, then the candidate
// recognizes the leader as legitimate and returns to follower
// state. If the term in the RPC is smaller than the candidate’s
// current term, then the candidate

package main

import (
	"context"
	"fmt"
)

func (n *Node) runCandidate(mainCtx context.Context, serverErrCh chan error) error {
	currentTerm := n.raftState.IncrementTerm()
	electionTimeout := n.raftState.NewElectionTimeout()

	childCtx, cancel := context.WithCancel(mainCtx)
	defer cancel()
	defer func() {
		n.logger.Info("candidate mode exiting")
	}()

	currentState := n.raftState.State()
	if currentState != StateCandidate {
		return nil
	}

	n.logger.Info("candidate state started successfully")
	_ = childCtx

	// TASK 1. Connect to all peers
	// TASK 2. Send a VoteRPC to all of them

	peers := connectToPeers(n.peers)
	if len(peers) == 0 {
		n.raftState.UpdateState(StateFollower)
		n.logger.Error("could not connect to any peers while a candidate")
		return nil
	}

	if candidate, granted := n.raftState.HasVotedFor(currentTerm); granted {
		panicMsg := fmt.Sprintf(
			`already granted vote for currentTerm: %d while in Candidate state
       Candidate is supposed to start election and vote for a higher/new term
       ---
       GrantedTo: %s,
       Raft diagnostics:
       %s
      `, currentTerm, candidate, n.raftState.String())

		panic(panicMsg)
	}

	n.raftState.GrantVoteTo(currentTerm, NodeId(n.id))
	timeoutCtx, cancel := context.WithTimeout(mainCtx, electionTimeout)
	defer cancel()

	wonElection := make(chan bool, 1)
	previousLogEntry := n.logStore.PreviousEntry()
	req := VoteRequest{
		Id:               NodeId(n.id),
		Term:             currentTerm,
		Result:           VoteResultVoteRequest,
		PreviousLogIndex: previousLogEntry.Index,
		PreviousLogTerm:  previousLogEntry.Term,
	}

	go collectOtherVotes(timeoutCtx, req, peers, wonElection)

	for {
		select {
		case <-mainCtx.Done():
			return mainCtx.Err()
		case <-timeoutCtx.Done():
			panic("election ctx has been reached")
		case err := <-serverErrCh:
			return err
		case result := <-wonElection:
			if result {
				n.logger.Info("candidate won election for new term, ascending to leader")
				n.raftState.UpdateState(StateLeader)
			} else {
				n.logger.Info("candidate lost election for new term going back to follower")
				n.raftState.UpdateState(StateFollower)
			}
		case payload := <-n.networkCh:
			// if we recv an appendEntry from a higher term
			_ = payload
			panic("candidate does not have handlers for payloads yet")
		}
	}
}

func connectToPeers(addrs []string) []*RPCPeer {
	panic("connectToPeers not yet implemented")
}

func collectOtherVotes(
	timeoutCtx context.Context,
	req VoteRequest,
	peers []*RPCPeer,
	won chan bool,
) {
	// spawn goroutines that collect the votes
	panic("vote collection not yet implemented")
}
