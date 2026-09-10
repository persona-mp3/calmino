package main

// How does a node, who has being restarted join the cluster. There are a couple
// of angles we can attack this at:
// 1. A new node that wasn't in the original cluster
// 2. A previous node that died for some reason, and got restarted
//
// First question, is it worth tacking the first scenario? It's not clear to me 
// that it would be easy, infact neither of them are easy. Because
// 1. The whole cluster now needs to know that this node has joined the cluster,
// we're talking about live reconfiguration here guys!
// 
//
// Second one seems more easier, because the others in the cluster don't need to 
// know or reconfigure, but the problem is 
// 1. How does a Leader know this Node genuinely just crashed and needs to be connected 
// again, because the first thing a restarted node does is wait for HB, and it's 
// not going to receive any if the leader doesn't know it exists
//
// DISCLAIMER: We're not persisting raft logs, even about crash data or stuff.
// Everything is volatile if not stored into the database,  but eventually, when
// we do, this problem might be easier
//
// One way  is to compare logs and terms right? If the crashed nodes prev logs
// matches 0, and if the term for that is also 0 and term is 1? then booya, the 
// leader adopts them and Updates their logs.
// But now, we're relying on that sentinel?? value or property to determine that.
// IF, we add log persistence, that will totally be invalidated. Another benefit
// it gives us is that, we can write the previous state of the node to disk and 
// read from there. If we even managed to write 'was follower, latest term: 19' 
// before our computer says bye-bye, will that help?
//
//
// Technically, I think that can hold? But now the leader now needs to know if 
// it already has this node on payroll, assigned to a worker before assigning it 
// to a worker. So the problem is 
// 1. How do we identify this Node is a lost cause, should there be any special things? 
// 2. How does the leader keep track of Followers that are gone??
// pain and suffering!
//
// I mean, per the protocol, if you reply with a higher term, the requester has 
// to concede, (including other checks) so by default they're just waiting for you
// to send heartbeats. It's also safe to say , that due to the prev design of 
// per-worker-node, that node worker would have already been demolished. so 
// WE don't need to check if we're already servicing them, we just add them. 
// There are prob some downsides we cant see yet, but this seems very viable
//
// - So if a node dies, the worker will also exit, that 'should' be gauranteed
// - When that node comes back up, it has amnesia, and will start sweating 
// - Problem is that, all nodes in the cluster will send a higher term, so who is the leader?
// there's no mention in the paper, that the leaderID is shared amongst each Follower?
// KKKKKKFCCCCCCCCUUUUUUUfFF
//
//
//
// Unless we just want to persist? who our last leader was? still back to the same 
// problem.
