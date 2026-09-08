package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// What if we just had a single loop sending out the heartbeats? instead of
// seperate workers?

func wokerLoop(
	leaderCtx context.Context,
	id NodeId,
	replicateCh chan replicate,
	peers []RPCConn,
	store LogStore,
	term uint64,
	raft *RaftState,
) {
	ticker := time.NewTicker(time.Millisecond * 300)
	defer func() {
		ticker.Stop()
		raft.UpdateState(StateFollower)
	}()

	// exit is closed  by sendHB when a majority of the cluster has ignored this node
	exit := make(chan struct{})

	prevLogEntry := store.PreviousEntry()
	commitIdx := store.CommitIndex()

	for {
		select {
		case <-leaderCtx.Done():
			return
		case replica := <-replicateCh:
			for _, peer := range peers {
				go func(success chan struct{}) {
					// TODO: need a timeout here incase the client is blocking to avoid leaks
					reply := AppendEntryReply{}
					if err := peer.Call("Server.AppendEntry", replica.req, &reply); err != nil {
						log.Println(err)
					}
					switch reply.Result {
					case RaftResultAcked, RaftResultLogsOutOfSync:
						select {
						case success <- struct{}{}:
						case <-time.After(SendChanTimeout):
							log.Println("dropping result from replication, receiver blocked")
						}
					}
				}(replica.replicated)
			}

			// update the values so ticker can read them
			prevLogEntry = store.PreviousEntry()
			commitIdx = store.CommitIndex()

		case <-ticker.C:
			req := AppendEntryRequest{
				Id:               id,
				Term:             term,
				PreviousLogIndex: prevLogEntry.Index,
				PreviousLogTerm:  prevLogEntry.Term,
				CommitIndex:      commitIdx,
			}
			// TODO: might need a goroutine pool here?
			go func() {
				err := sendHB(leaderCtx, peers, req, exit)
				if err != nil {
					log.Println("[error] from sendHB", err)
					return
				}
			}()
			ticker.Reset(time.Millisecond * 300)
		case <-exit:
			log.Println("[info] exit channel fired")
			return
		}
	}
}

func sendHB(leaderCtx context.Context, peers []RPCConn, req AppendEntryRequest, exit chan struct{}) error {
	failed := atomic.Uint64{}
	wg := sync.WaitGroup{}
	for _, peer := range peers {
		wg.Add(1)
		go func(peer RPCConn, failed *atomic.Uint64) {
			reply := AppendEntryReply{}
			if err := peer.Call("Server.AppendEntry", req, &reply); err != nil {
				log.Println("[error] sending heartbeat", err)
				failed.Add(1)
				return
			}

			switch reply.Result {
			case RaftResultAcked, RaftResultLogsOutOfSync:
			default:
				log.Printf("[info] was not acked by follower: %+v\n", reply)
				failed.Add(1)
			}
		}(peer, &failed)
	}

	done := make(chan struct{}, 1)
	go func() {
		wg.Wait()
		close(done)
		log.Println("all background routines returned")
	}()

	select {
	case <-leaderCtx.Done():
		return fmt.Errorf("from sendHB: %w", leaderCtx.Err())
	case <-done:
	}
	ackedByMajority := failed.Load() <= uint64(len(peers))
	if !ackedByMajority {
		close(exit)
	}
	return nil
}
