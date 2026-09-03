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
	"log"
	"log/slog"
	"sync/atomic"
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

	// TASK 3: Write handlers for candidate
	// TASK 4. Impl the VoteRPC Method

	rpcPeers := connectToPeers("tcp", n.peers)
	if len(rpcPeers) == 0 {
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

	n.logger.Info("parading for", slog.Uint64("term", currentTerm))
	go collectOtherVotes(timeoutCtx, req, len(n.peers), rpcPeers, wonElection)
	handler := NewCandidateHandler(NodeId(n.id), n.logger)

	for {
		select {
		case <-mainCtx.Done():
			return mainCtx.Err()
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
			return nil
		case payload := <-n.networkCh:
			// if we recv an appendEntry from a higher term
			_ = payload
			var response RPCReply

			switch req := payload.payload.(type) {
			case AppendEntryRequest:
				reply, err := handler.AppendEntry(&req, n.raftState, n.logStore)
				if err != nil {
					panic(err)
				}
				if reply.Result == RaftResultAcked {
					n.raftState.UpdateState(StateFollower)
					payload.reply <- response
					n.logger.Info("dropping down to follower from candidate")
					return nil
				}
				response = RPCReply{kind: RPCKindAppendEntry, payload: reply}
			case VoteRequest:
				reply, err := handler.Vote(&req, n.raftState, n.logStore)
				if err != nil {
					panic(err)
				}
				response = RPCReply{kind: RPCKindVote, payload: reply}
				if reply.Result == VoteResultVoteGranted {
					payload.reply <- response
					n.logger.Info("dropping down to follower from candidate")
					return nil
				}

			case SnapshotRequest:
				reply, err := handler.Snapshot(&req, n.raftState, n.logStore)
				if err != nil {
					panic(err)
				}
				response = RPCReply{kind: RPCKindSnapshot, payload: reply}
			default:
				panicMsg := fmt.Sprintf(
					"recvd unsupported or unimplemented payload as candidate:\n%+v\n", req,
				)
				panic(panicMsg)

			}
			payload.reply <- response

		}
	}
}

func routePayload(
	payload RPCPayload,
	raftState *RaftState,
	logStore LogStore,
	handler Handler,
) (RPCReply, error) {

	switch req := payload.payload.(type) {
	case AppendEntryRequest:
		reply, err := handler.AppendEntry(&req, raftState, logStore)
		return RPCReply{kind: RPCKindAppendEntry, payload: reply}, err
	case VoteRequest:
		reply, err := handler.Vote(&req, raftState, logStore)
		return RPCReply{kind: RPCKindVote, payload: reply}, err
	case SnapshotRequest:
		reply, err := handler.Snapshot(&req, raftState, logStore)
		return RPCReply{kind: RPCKindSnapshot, payload: reply}, err
	default:
		panicMsg := fmt.Sprintf(
			"recvd unsupported or unimplemented payload as candidate:\n%+v\n", req,
		)
		panic(panicMsg)
	}
}

func collectOtherVotes(
	timeoutCtx context.Context,
	req VoteRequest,
	clusterSize int,
	peers []*RPCPeer,
	won chan bool,
) {
	collectedVotes := atomic.Uint64{}
	collectedVotes.Add(1)
	newVote := make(chan struct{}, len(peers)*2)
	majorityVote := (clusterSize / 2) + 1

	for _, peer := range peers {
		go func(peer *RPCPeer, voteCh chan struct{}) {
			reply := VoteReply{}
			if err := peer.Call("Server.VoteRPC", req, &reply); err != nil {
				log.Println("[error] failed to call VoteRPC", err)
				return
			}

			if reply.Result != VoteResultVoteGranted {
				log.Println("[warn] vote request was not granted", reply)
				return
			}
			voteCh <- struct{}{}
			log.Println("[debug] sent vote to channel", reply)
		}(peer, newVote)
	}

	for {
		if collectedVotes.Load() >= uint64(majorityVote) {
			won <- true
			return
		} else if timeoutErr := timeoutCtx.Err(); timeoutErr != nil {
			log.Println("[debug] timeoutCtx done:", timeoutErr)
			break
		}

		select {
		case <-newVote:
			collectedVotes.Add(1)
		case <-timeoutCtx.Done():
			log.Println("[debug] timeoutCtx fired")
		}
	}

	electionWon := collectedVotes.Load() >= uint64(majorityVote)
	won <- electionWon
}
