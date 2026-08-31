package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

func RunCluster(rc *RawConfig) error {
	if rc.Mode == singleProcess {
		singleProcessCluster(rc)
		return nil
	}
	return fmt.Errorf("%s not yet implemented. please use single_process", string(rc.Mode))
}

func singleProcessCluster(rc *RawConfig) {
	idPeers := buildIndexedPeerMap(rc.Addrs)
	allNodes := []*Node{}

	for idx, addrPeers := range idPeers {
		for addr, peers := range addrPeers {
			id := fmt.Sprintf("%d", idx+1)
			var out io.Writer
			if rc.Persist {
				out = createFileWithName(fmt.Sprintf("log-file-%s", id))
			} else {
				out = os.Stdout
			}
			logStore := NewLogStore()
			d := randomDuration(time.Millisecond)
			raftState := NewRaftState(d)
			config, err := rc.ToConfig()
			if err != nil {
				log.Panic("error creating config:", err)
			}
			config.Out = out
			node := NewNode(id, addr, peers, &logStore, raftState, *config)
			allNodes = append(allNodes, node)
		}
	}

	nodeWg := sync.WaitGroup{}
	for _, node := range allNodes {
		nodeWg.Go(func() {
			if err := node.Run(); err != nil {
				log.Println("could not start node:", node.id, err)
			}
		})
	}

	nodeWg.Wait()
	log.Println("all nodes returned")
}

func singleNodeCluster(rc *RawConfig) {}

func createFileWithName(name string) io.Writer {
	f, err := os.Create(name)
	if err != nil {
		f = os.Stdout
		log.Println(
			"[warn] could not create file: %s, reason: %s. using stdout",
			name,
			err)
	}

	return f
}
