package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func RunCluster(rc *RawConfig) error {
	go func() {
		log.Println("[pprof] starting pprof server at http://", rc.HTTPPprofAddr)
		err := startPprofServer(rc.HTTPPprofAddr)
		if err != nil {
			log.Println("[error] ", err)
		}
	}()
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
			node := NewNode(id, addr, peers, logStore, raftState, *config)
			allNodes = append(allNodes, node)
		}
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGKILL)
	defer cancel()

	nodeWg := sync.WaitGroup{}
	for _, node := range allNodes {
		nodeWg.Go(func() {
			if err := node.Start(ctx); err != nil {
				log.Println("could not start node:", node.id, err)
			}
		})
	}

	nodeWg.Wait()
	log.Println("all nodes returned")
}

func multiProcessCluster(rc *RawConfig) {
	// the multiprocess cluster needs all nodes to run on different processes.
	// We still want this go process to be able to control the others, so we
	// might need to get a bit creative here.
	// For each node, we'll have to do something like this
	// os.Exec("./calmino").args("--config", "node-1-config.toml")
	// the config can be generated at runtime. The need for the config is just to
	// tell calmino which ones it should run? or intead of a whole new config we
	// just run
	// os.Exec("./calmino").args("--id", "1")
	// where id is the position in the addr.
	// So when we recv --id we just know to run a single and at at the given index.
	// The single one can also be used to run a single node so, we end up with best
	// of both worlds.
	//
	// The next thing is
	// 1. process control,
	// 2. http_pproff_addr [we can enforce the config schema to provide addrs for other nodes]

	// Process control might be tricky bcus I haven't done it in Go before where
	// a user doesn't need to kill all other processes directly, to stop the whole
	// cluster. Instead if quit this program, that kills all other processes
	// forked from it OR that it spawned.
	// This also buys room for a control plane, where we can decide to kill a specific
	// node in the cluster. WOW there's so much crazy things that can be done here
	// I can almost see k8 control plane staring at me from the corner. The peak of
	// this would be able to design the control plane across the network instead.
	// Where I can kill and restart nodes

}
