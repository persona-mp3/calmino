package main

import (
	"context"
	"time"
)

// QUESTION
// This is supposed to be an extension of the leader, where they're responsible
// for sending heartbeat to all connected nodes in the cluster. Each worker should
// also be able to assit the leader in replicating data across nodes when client
// commands come in. But how would a worker know about a change in commit index?
// Do we just give them a a pointer value to share? Because raftState.CommitIdx()
// returns the latest, but theres a mutex there, and calling that amongs other 
// routines at random timeouts is expensive
//
// type leader struct {
// 		handler        Handler
// 		network        chan RPCNetwork
// 		raft           *RaftState
// 		store          LogStore
// 		connectedPeers []RPCConn
// 		workers        []*Worker
// }
//
// func (l *leader) run(ctx context.Context) {
// 	for {
// 		select {
// 			case <-ctx.Done():
// 			 	return
// 			 case payload := <-l.network:
// 				switch req := payload.payload.(type) {
// 					case Command:
// 						log := l.store.ReqToLog(req)
// 						idx := l.store.Append(log)
// 						log.Index = idx
// 						success := startReplication(500 * time.Millisecond, log, l.workers)
// 						if !success {
// 							payload.reply <- RPCReply{Command: {"error, please try again"}}
// 							continue
// 						}
//
// 					result, commitIdx := l.store.Flush()
// 					payload.reply <- RPCReply{Command: {result}}
// 				default:
// 					resp := handler(payload)
// 					payload.reply <- resp
// 				}
// 		}
// 	}
// }
//
// // not sure yet, might want startReplication to take a request instead
// func startReplication(ctx context.Context, dur time.Duration, log Log, workers []*Worker) bool {
// 	 done := make(chan struct{}, len(workers)*2) // make sure others can still send
// 	 replicateCmd := replicate{log: log, done: done}
// 	 sendTicker := ticker.NewTicker(SendChannelTimeout)
// 	 for _, worker := range workers {
// 	 	sendTicker.Reset(SendChannelTimeout)
// 	 	go func(){
// 	 		select {
// 	 		case worker.replicate <- replicateCmd:
// 	 		case <-sendTicker.C:
// 	 			// would not want to use ctx here because we dont need the cancel and 
// 	 			// will need to create a new one for each iteration. ticker automatically 
// 	 			// resolves both or a timer instead?
// 	 			logger.Warn("worker has been blocked, dropping replicate command")
//
// 	 		}
// 	 	}()
// 	 }
//
// 	 replicationTimeoutCtx, cancel := context.WithTimeout(ctx, dur)
// 	 defer cancel()
//
// 	 votes := 1
// 	 for {
// 	 	 if err := replicationTimeoutCtx.Error() != nil {
// 	 	 	break
// 	 	 }
// 	 	 if votes >= expectedMajority {
// 	 	 	return true
// 	 	 }
// 	 	 select {
// 	 	 case <-done:
// 	 	 	votes += 1
// 	 	 case <-replicationTimeoutCtx:
// 	 	 	break
// 	 	 }
// 	 }
//
// 	 return votes >= expectedMajority
// }
//
// func (l *leader) startWorkers(
// 	ctx context.Context, 
// 	commitTracker *atomic.Uint64, 
// 	peer []RPCConn,
// 	newEntry chan replicate,
// 	) {
// 		worker := newWorker(commitTracker, l.raft, l.store)
// 		workerWg := sync.WaitGroup{}
// 		for _, peer := range peer {
// 		 	wg.Go(func(){ worker.Run(ctx) }) 
// 		}
//
// 		wg.Wait()
// 		logger.Debug("all workers returned")
// }

type replicate struct {
	done chan struct{}
	req  AppendEntryRequest
}

type Worker struct {
	term     uint64
	peer     RPCConn
	logStore LogStore
}

func (w *Worker) run(ctx context.Context, heartbeat time.Duration) {
	ticker := time.NewTicker(heartbeat)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			panic("worker.run for hearbeat ticker not implemented yet")
		}
	}
}
