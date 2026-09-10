package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func RunCluster(mainCtx context.Context, rc *RawConfig) error {
	if len(rc.HTTPPprofAddr) == 0 {
		return fmt.Errorf("http_pprof_addr cannot be empty, please provide an addr")
	}

	switch rc.Mode {
	case clusterModeSingleProcess:
		go func() {
			log.Println("[pprof] starting pprof server at http://", rc.HTTPPprofAddr)
			err := startPprofServer(rc.HTTPPprofAddr[0])
			if err != nil {
				log.Println("[error] ", err)
			}
		}()
		singleProcessCluster(rc)
		return nil
	case clusterModeMultiProcess:
		runMultiProcessCluster(mainCtx, rc)
		return nil
	}
	return fmt.Errorf("%s not yet implemented. please use single_process", string(rc.Mode))
}

func singleProcessCluster(rc *RawConfig) {
	allNodes := createNodes(rc)
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

func runMultiProcessCluster(ctx context.Context, rc *RawConfig) {
	clusterSize := len(rc.Addrs)

	allCommands := []*exec.Cmd{}
	cmdWg := sync.WaitGroup{}
	for idx := range clusterSize {
		cmd := exec.CommandContext(ctx, "./calmino", "--nodeId", strconv.Itoa(idx+1))
		allCommands = append(allCommands, cmd)
		cmdWg.Go(func() {
			if err := cmd.Run(); err != nil {
				log.Printf("[error] starting %d. reason: %s\n", idx, err)
			}
		})
	}

	fmt.Println("waiting for all nodes to complete")
	cmdWg.Wait()
	fmt.Println("all nodes completed")
}

func createNodes(rc *RawConfig) []*Node {

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
				log.Fatal("error creating config:", err)
			}
			config.Out = out
			node := NewNode(id, addr, peers, logStore, raftState, *config)
			allNodes = append(allNodes, node)
		}
	}

	return allNodes
}
